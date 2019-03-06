// +build linux

package collector

import (
	"io/ioutil"
	"os"
	"strings"
	"regexp"
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
func (linux *Linux) GetLoginLog() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetProcessList get process list
func (linux *Linux) GetProcessList() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
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
