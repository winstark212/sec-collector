package main

import (
	"fmt"
	"github.com/winstark212/sec-collector/collector"
)


func main() {
	collect := collector.Windows{}
	fmt.Print(collect.GetComInfo())
}