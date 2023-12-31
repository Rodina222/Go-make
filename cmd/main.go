package main

import (
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-Rodina/internal"
)

const CommandFormat = "gomake -f Makefile -t target"

func main() {

	// parse command line
	target, filePath, err := internal.ParseCommandLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse command line: %v\n", fmt.Errorf("%w", err))
		fmt.Println(CommandFormat)
		os.Exit(1)
	}

	// create new graph
	graph := internal.NewGraph()

	// call execute method of the graph
	err = graph.Execute(filePath, target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
