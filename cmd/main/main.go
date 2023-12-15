package main

import (
	"flag"

	"github.com/Yeuoly/shukuangkuang/internal"
)

func parseArgs() internal.ShukuangkuangArgs {
	args := internal.ShukuangkuangArgs{}
	flag.BoolVar(&args.LogicCoreMode, "mode", true, "use logic cpu cores mode")
	flag.Parse()
	return args
}

func main() {
	args := parseArgs()
	core := internal.Shukuangkuang{}
	core.Run(args)
}
