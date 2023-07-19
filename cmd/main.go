package main

import (
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-Rodina/internal"
)

const CommandFormat = "./gomake -f Makefile -t target"

func main() {

	target, filePath, err := internal.ParseCommandLine()

	graph := internal.NewGraph()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse command line: %v\n", err)
		fmt.Println(CommandFormat)
		os.Exit(1)
	}

	// parse the Makefile (create the graph)
	err = graph.ParseMakeFile(filePath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check cyclic dependency
	err = graph.CheckCyclicDependency()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//check all targets have commands
	err = graph.CheckCmds()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// get the order of executing commands and execute them
	err = graph.ExecuteInOrder(target)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
