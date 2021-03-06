// +build linux

package collector

import (
	"io/ioutil"
	"os"
	"time"
	"fmt"
	"log"
	"strings"
	"strconv"
	"regexp"
	"bytes"
	"encoding/binary"
	"github.com/winstark212/sec-collector/common"
)

// Linux struct
type Linux struct {
}

// GetComInfo get computer information
func (linux *Linux) GetComInfo() (info ComputerInfo) {
	info.IP = LocalIP
	info.Hostname, _ = os.Hostname()
	out := common.Cmdexec("uname -r")
	dat, err := ioutil.ReadFile("/etc/redhat-release")
	if err != nil {
		dat, _ = ioutil.ReadFile("/etc/issue")
		issue := strings.SplitN(string(dat), "\n", 2)[0]
		out2 := common.Cmdexec("uname -m")
		info.System = issue + " " + out + out2
	} else {
		info.System = string(dat) + " " + out
	}
	// discern(&info)
	return info
}

// GetCrontab Get scheduled tasks
func (linux *Linux) GetCrontab() (resultData []map[string]string) {
	//system planing task
	dat, err := ioutil.ReadFile("/etc/crontab")
	if err != nil {
		return resultData
	}
	cronList := strings.Split(string(dat), "\n")
	for _, info := range cronList {
		if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 6 {
			continue
		}
		s := strings.SplitN(info, " ", 7)
		rule := strings.Split(info, " "+s[5])[0]
		m := map[string]string{"command": s[6], "user": s[5], "rule": rule}
		resultData = append(resultData, m)
	}

	// user scheduled task
	dir, err := ioutil.ReadDir("/var/spool/cron/")
	if err != nil {
		return resultData
	}
	for _, f := range dir {
		if f.IsDir() {
			continue
		}
		dat, err = ioutil.ReadFile("/var/spool/cron/" + f.Name())
		if err != nil {
			continue
		}
		cronList = strings.Split(string(dat), "\n")
		for _, info := range cronList {
			if strings.HasPrefix(info, "#") || strings.Count(info, " ") < 5 {
				continue
			}
			s := strings.SplitN(info, " ", 6)
			rule := strings.Split(info, " "+s[5])[0]
			m := map[string]string{"command": s[5], "user": f.Name(), "rule": rule}
			resultData = append(resultData, m)
		}
	}
	return resultData
}

// GetListening Get tcp listening port
func (linux *Linux) GetListening() (resultData []map[string]string) {
	listeningStr := common.Cmdexec("ss -nltp")
	listeningList := strings.Split(listeningStr, "\n")
	if len(listeningList) < 2 {
		return
	}
	for _, info := range listeningList[1 : len(listeningList)-1] {
		if strings.Contains(info, "127.0.0.1") {
			continue
		}
		m := make(map[string]string)
		reg := regexp.MustCompile("\\s+")
		info = reg.ReplaceAllString(strings.TrimSpace(info), " ")
		s := strings.Split(info, " ")
		if len(s) < 6 {
			continue
		}
		m["proto"] = "TCP"
		if strings.Contains(s[3],"::"){
			m["address"] = strings.Replace(s[3], "::", "0.0.0.0", 1)
		}else{
			m["address"] = strings.Replace(s[3], "*", "0.0.0.0", 1)
		}
		b := false
		for _,v:= range resultData{
			if v["address"] == m["address"]{
				b = true
				break
			}
		}
		if b{
			continue
		}
		reg = regexp.MustCompile(`users:\(\("(.*?)",(.*?),.*?\)`)
		r := reg.FindSubmatch([]byte(s[5]))
		if strings.Contains(string(r[2]), "=") {
			m["pid"] = strings.SplitN(string(r[2]), "=", 2)[1]
		} else {
			m["pid"] = string(r[2])
		}
		m["name"] = string(r[1])
		resultData = append(resultData, m)
	}
	return resultData
}

// GetLoginLog get login log
func (linux *Linux) GetLoginLog() (resultData []map[string]string) {
	resultData = getLast(Config.Lasttime)
	resultData = append(resultData, getLastb(Config.Lasttime)...)
	return
}

// GetProcessList get process list
func (linux *Linux) GetProcessList() (resultData []map[string]string) {
	var dirs []string
	var err error
	dirs, err = dirsUnder("/proc")
	if err != nil || len(dirs) == 0 {
		return
	}
	for _, v := range dirs {
		pid, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		statusInfo := getStatus(pid)
		command := getcmdline(pid)
		m := make(map[string]string)
		m["pid"] = v
		m["ppid"] = statusInfo["PPid"]
		m["name"] = statusInfo["Name"]
		m["command"] = command
		resultData = append(resultData, m)
	}
	return
}

// GetServiceInfo get service list
func (linux *Linux) GetServiceInfo() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetStartup get startup service
func (linux *Linux) GetStartup() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetUser get user list
func (linux *Linux) GetUser() (resultData []map[string]string) {
	data, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return resultData
	}
	userList := strings.Split(string(data), "\n")
	if len(userList) < 2 {
		return
	}
	for _, info := range userList[0 : len(userList)-2] {
		if strings.Contains(info, "/nologin") {
			continue
		}
		s := strings.SplitN(info, ":", 2)
		m := map[string]string{"name": s[0], "description": s[1]}
		resultData = append(resultData, m)
	}
	return resultData
}


// additional functions

func getLast(t string) (result []map[string]string) {
	var timestamp int64
	if t == "all" {
		timestamp = 615147123
	} else {
		ti, _ := time.Parse("2006-01-02T15:04:05Z07:00", t)
		timestamp = ti.Unix()
	}
	wtmpFile, err := os.Open("/var/log/wtmp")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer wtmpFile.Close()
	for {
		wtmp := new(utmp)
		err = binary.Read(wtmpFile, binary.LittleEndian, wtmp)
		if err != nil {
			break
		}
		if wtmp.UtType == 7 && int64(wtmp.UtTv.TvSec) > timestamp {
			m := make(map[string]string)
			m["status"] = "true"
			m["remote"] = string(bytes.TrimRight(wtmp.UtHost[:], "\x00"))
			if m["remote"] == "" {
				continue
			}
			m["username"] = string(bytes.TrimRight(wtmp.UtUser[:], "\x00"))
			m["time"] = time.Unix(int64(wtmp.UtTv.TvSec), 0).Format("2006-01-02T15:04:05Z07:00")
			result = append(result, m)
		}
	}
	return result
}

func getLastb(t string) (result []map[string]string) {
	var cmd string
	ti, _ := time.Parse("2006-01-02T15:04:05Z07:00", t)
	if t == "all" {
		cmd = "lastb --time-format iso"
	} else {
		cmd = fmt.Sprintf("lastb -s %s --time-format iso", ti.Format("20060102150405"))
	}
	out := common.Cmdexec(cmd)
	logList := strings.Split(out, "\n")
	for _, v := range logList[0 : len(logList)-3] {
		m := make(map[string]string)
		reg := regexp.MustCompile("\\s+")
		v = reg.ReplaceAllString(strings.TrimSpace(v), " ")
		s := strings.Split(v, " ")
		if len(s) < 4 {
			continue
		}
		m["status"] = "false"
		m["username"] = s[0]
		m["remote"] = s[2]
		t, _ := time.Parse("2006-01-02T15:04:05Z0700", s[3])
		m["time"] = t.Format("2006-01-02T15:04:05Z07:00")
		result = append(result, m)
	}
	return
}

func getcmdline(pid int) string {
	cmdlineFile := fmt.Sprintf("/proc/%d/cmdline", pid)
	cmdlineBytes, e := ioutil.ReadFile(cmdlineFile)
	if e != nil {
		return ""
	}
	cmdlineBytesLen := len(cmdlineBytes)
	if cmdlineBytesLen == 0 {
		return ""
	}
	for i, v := range cmdlineBytes {
		if v == 0 {
			cmdlineBytes[i] = 0x20
		}
	}
	return strings.TrimSpace(string(cmdlineBytes))
}
func getStatus(pid int) (status map[string]string) {
	status = make(map[string]string)
	statusFile := fmt.Sprintf("/proc/%d/status", pid)
	var content []byte
	var err error
	content, err = ioutil.ReadFile(statusFile)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(content), "\n") {
		if strings.Contains(line, ":") {
			kv := strings.SplitN(line, ":", 2)
			status[kv[0]] = strings.TrimSpace(kv[1])
		}
	}
	return
}

func dirsUnder(dirPath string) ([]string, error) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}
	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}
	return ret, nil
}
