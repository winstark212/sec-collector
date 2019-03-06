// +build linux

package collect

import (
	"io/ioutil"
	"os"
	"strings"
	"regexp"
	"github.com/winstark212/sec-collector/common"
)

type utmp struct {
	UtType uint32
	UtPid  uint32    // PID of login process
	UtLine [32]byte  // device name of tty - "/dev/"
	UtID   [4]byte   // init id or abbrev. ttyname
	UtUser [32]byte  // user name
	UtHost [256]byte // hostname for remote login
	UtExit struct {
		ETermination uint16 // process termination status
		EExit        uint16 // process exit status
	}
	UtSession uint32 // Session ID, used for windowing
	UtTv      struct {
		TvSec  uint32 /* Seconds */
		TvUsec uint32 /* Microseconds */
	}
	UtAddrV6 [4]uint32 // IP address of remote host
	Unused   [20]byte  // Reserved for future use
}

// GetComInfo get computer information
func GetComInfo() (info common.ComputerInfo) {
	info.IP = common.LocalIP
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
func GetCrontab() (resultData []map[string]string) {
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
func GetListening() (resultData []map[string]string) {
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



