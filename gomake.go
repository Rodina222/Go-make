package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/codescalersinternships/gomake-Rodina/internal"
)

var (
	ErrTargetNotFound = errors.New("target is not found to be executed")
	ErrParseMakeFile  = errors.New("failed to parse the makefile")
)

const CommandFormat = "./gomake -f Makefile -t target"

func main() {

	graph := internal.NewGraph()
	target, filePath, err := ParseCommandLine()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse command line: %v\n", err)
		fmt.Println(CommandFormat)
		return
	}

	// parse the Makefile (create the graph)
	err = graph.ParseMakeFile(filePath)

	if err != nil {
		wrappedErr := fmt.Errorf("%w line %d", ErrParseMakeFile, err)
		fmt.Println(wrappedErr)
		return
	}

	// check the input target is available in the makefile
	err = graph.CheckTarget(target)

	if err != nil {
		fmt.Println(err)
		return

	}

	// check cyclic dependency
	graph.CheckCyclicDependency()

	// get the order of executing commands
	graph.OrderOfExecution(target)

	//execute the target and its dependencies (if found)
	internal.ExecCommand(target)

}

// ParseCommandLine parses the command line to extract the filename and target
func ParseCommandLine() (string, string, error) {

	filePath := flag.String("f", "", "file path to be opened")
	target := flag.String("t", "", "target to be executed")
	flag.Parse()

	// check that there is a target
	if *target == "" {
		return "", "", ErrTargetNotFound
	}

	// specify the makefile if it is not given
	if *filePath == "" {
		*filePath = "Makefile"
	}

	return *target, *filePath, nil

}
