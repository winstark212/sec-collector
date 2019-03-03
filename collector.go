package main

import (
	"fmt"
	"github.com/winstark212/sec-collector/common"
)


func main() {
	result := common.Cmdexec("x")
	fmt.Printf(result)
}