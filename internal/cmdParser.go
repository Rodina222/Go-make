package internal

import (
	"flag"
)

// ParseCommandLine parses the command line to extract the filename and target
func ParseCommandLine() (string, string, error) {

	var filePath, target string

	flag.StringVar(&filePath, "f", "Makefile", "file path to be opened")
	flag.StringVar(&target, "t", "", "target to be executed")

	flag.Parse()

	// check that there is a target
	if target == "" {
		return "", "", ErrTargetNotFound
	}

	return target, filePath, nil

}
