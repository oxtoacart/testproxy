package main

import (
	_ "github.com/oxtoacart/testproxy/proxy"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(4)

	time.Sleep(10000 * time.Hour)
}
