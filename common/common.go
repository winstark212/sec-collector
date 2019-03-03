package common

import (
	"os/exec"
	// "runtime"
	"log"
	"strings"
)

// ComputerInfo Computer information structure
type ComputerInfo struct {
	IP       string   // IP address
	System   string   // Operating system
	Hostname string   // Computer name
	Type     string   // Server type
	Path     []string // Web directory
}

func Cmdexec(cmd string) string {
	var c *exec.Cmd
	var data string
	// system := runtime.GOOS
	argArray := strings.Split(cmd, " ")
	c = exec.Command(argArray[0], argArray[1:]...)
	out, err := c.CombinedOutput()
	if err != nil {
		log.Print(err)
	}
	data = string(out)
	return data
}
