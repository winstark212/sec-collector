package main

import (
	"fmt"
	"github.com/winstark212/sec-collector/collect"
)


func main() {
	fmt.Print(collect.GetCrontab())
}