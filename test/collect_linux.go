package main

import (
	"time"
	"log"
	// "encoding/binary"
	// "os"
	// "bytes"
	"github.com/winstark212/sec-collector/common"
	"strings"
	"fmt"
	"regexp"
	"encoding/json"

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

// additional functions
func getLastb(t string) (result []map[string]string) {
	var cmd string
	ti, _ := time.Parse("2006-01-02T15:04:05Z07:00", t)
	if t == "all" {
		cmd = "lastb --time-format iso"
	} else {
		cmd = fmt.Sprintf("lastb -s %s --time-format iso", ti.Format("20060102150405"))
	}
	out := common.Cmdexec(cmd)
	fmt.Print("Result: " + out)
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
func main() {
	result := getLastb("all")
	jsonResultData, _ := json.Marshal(result)
	log.Print(string(jsonResultData))
}