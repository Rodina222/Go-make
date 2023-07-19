package main

import (
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-Rodina/internal"
)

const CommandFormat = "gomake -f Makefile -t target"

func main() {

	target, filePath, err := internal.ParseCommandLine()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse command line: %v\n", fmt.Errorf("%w", err))
		fmt.Println(CommandFormat)
		os.Exit(1)
	}

	graph := internal.NewGraph()

	err = graph.Execute(filePath, target)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
