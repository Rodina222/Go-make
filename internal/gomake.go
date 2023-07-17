package internal

import (
	"flag"
)

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
