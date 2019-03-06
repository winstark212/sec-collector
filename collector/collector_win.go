package collector

import (
 "fmt"
)
// Windows struct
type Windows struct {
}

// GetComInfo get computer information
func (windows *Windows) GetComInfo() (info ComputerInfo) {
	fmt.Print("Get computer information")
	return info
}

// GetCrontab Get scheduled tasks
func (windows *Windows) GetCrontab() ( []map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetListening Get tcp listening port
func (windows *Windows) GetListening() ( []map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetLoginLog get login log
func (windows *Windows) GetLoginLog() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetProcessList get process list
func (windows *Windows) GetProcessList() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetServiceInfo get service list
func (windows *Windows) GetServiceInfo() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetStartup get startup service
func (windows *Windows) GetStartup() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}

// GetUser get user list
func (windows *Windows) GetUser() ([]map[string]string) {
	var resultData []map[string]string
	return resultData
}
