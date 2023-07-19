package internal

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecCommand executes a command line when it is called
func ExecCommand(command string) error {

	// check if the command is suppressed
	suppressedCmd := true
	if strings.HasPrefix(command, "@") {
		suppressedCmd = false
		command = strings.TrimPrefix(command, "@")
	}

	// split the command into words
	parts := strings.Fields(command)

	// get the name of the command
	cmdName := parts[0]

	// find the path to the command using LookPath
	cmdPath, err := exec.LookPath(cmdName)

	if err != nil {
		return err
	}

	// create a new Command object with the command and its arguments
	cmd := exec.Command(cmdPath, parts[1:]...)

	// set the output to print to the console
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// print the command if not suppressed
	if !suppressedCmd {
		fmt.Println(command)
	}

	// execute the command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
