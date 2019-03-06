package main

import (
	"fmt"
	"github.com/winstark212/sec-collector/common"
	"github.com/winstark212/sec-collector/test"
)


func main() {
	result := common.Cmdexec("x")
	fmt.Printf(result)
	fmt.Print(test.Hello())
}