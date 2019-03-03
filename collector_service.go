package main


type Collector interface {
	GetComputerInfo() map[string]string
	GetCrontab() map[string]string
}



