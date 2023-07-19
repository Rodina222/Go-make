package main

import (
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-Rodina/internal"
)

const CommandFormat = "./gomake -f Makefile -t target"

func main() {

	graph := internal.NewGraph()

	target, filePath, err := internal.ParseCommandLine()

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

	err = graph.Execute(target)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
