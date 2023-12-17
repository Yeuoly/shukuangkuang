package main

import (
	"flag"
	"fmt"

	"github.com/Yeuoly/shukuangkuang/internal"
)

func parseArgs() internal.ShukuangkuangArgs {
	args := internal.ShukuangkuangArgs{}
	flag.BoolVar(&args.LogicCoreMode, "mode", true, "use logic cpu cores mode")
	flag.BoolVar(&args.Help, "help", false, "print help")
	flag.Parse()
	return args
}

func printHelp() {
	fmt.Println("Welcome to use shukuangkuang!")
	fmt.Println("This is a tool to monitor your system status in terminal like the Task Manager in Windows.")
	fmt.Println("There is 2 modes in this tool: CPU mode and Memory mode.")
	fmt.Println("You can switch between them by pressing 'c' and 'm' on your keyboard.")
	fmt.Println("You can also press 'q' or 'Ctrl + c' to quit.")
	fmt.Println("Options:")
	flag.Usage()
}

func main() {
	args := parseArgs()
	if args.Help {
		printHelp()
		return
	}
	core := internal.Shukuangkuang{}
	core.Run(args)
}
