package main

import (
	"fmt"
	"time"

	"github.com/Yeuoly/shukuangkuang/internal"
)

func main() {
	memloader := &internal.MemoryStatusLoader{}
	memloader.Init()
	start := time.Now()
	usages, procs, _ := memloader.GetMemoryStatus()
	fmt.Printf("time: %v\n", time.Since(start))
	fmt.Printf("usages: %v\n", usages)

	max := 5
	if len(procs) < max {
		max = len(procs)
	}

	for i := 0; i < max; i++ {
		fmt.Printf("pid:%v\tname:%v\tusage:%v\t\n", procs[i].Pid, procs[i].Name, procs[i].MemoryUsage)
	}
}
