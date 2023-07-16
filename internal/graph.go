package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (

	// ErrTargetNotFound is returned when there is no input target or the input target doesn't exist
	ErrTargetNotFound = errors.New("target is not found to be executed")
)

// Vertex represents a target in a graph of vertices/targets that has dependencies and commands
type Vertex struct {

	// slice of dependencies
	dependencies []string

	// slice of commands
	cmds []string
}

// Graph represents a graph of vertices/targets
type Graph struct {
	vertices map[string]Vertex
}

// NewVertex returns a vertex
func NewVertex() Vertex {

	return Vertex{dependencies: make([]string, 0),
		cmds: make([]string, 0),
	}
}

// NewGraph returns a graph
func NewGraph() Graph {

	return Graph{
		vertices: make(map[string]Vertex),
	}

}

// AddCommand adds a command to a target in the graph
func (vertex *Vertex) AddCommand(command string) {

	vertex.cmds = append(vertex.cmds, command)
}

// AddDependencies adds the dependencies to a target in the graph
func (vertex *Vertex) AddDependencies(line string) {

	dependencies := strings.Split(line, " ")

	vertex.dependencies = append(vertex.dependencies, dependencies...)

}

// ParseMakeFile it parses the input file in the command line
func (graph *Graph) ParseMakeFile(filepath string) error {

	target := ""

	//opening the file
	content, err := os.Open(filepath)

	if err != nil {
		return errors.New("failed to open the file")
	}

	defer content.Close()

	// create a scanner to read the file line by line
	scanner := bufio.NewScanner(content)

	// create a vertex in the graph
	vertex := NewVertex()

	// scanning the file line by line
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// skipping empty lines and comments
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, ":") && line[0] != ' ' {

			// create a new vertex in the graph
			vertex = NewVertex()

			// specify the position of colon  ":"
			colonIndex := strings.Index(line, ":")

			//extract the target from the line
			target = strings.TrimSpace(line[:colonIndex])

			if target == "" {
				return ErrTargetNotFound
			}

			// add dependencies (if found)
			if len(line) > colonIndex+1 {

				deps := strings.TrimSpace(line[colonIndex+1:])
				vertex.AddDependencies(deps)

			}
			continue
		}

		vertex.AddCommand(line)

		// add the vertex to the graph
		graph.vertices[target] = vertex

	}

	return nil

}

// CheckCyclicDependency checks if a cyclic dependency exists or not and exits once it found one
func (graph *Graph) CheckCyclicDependency() {

	for target, vertex := range graph.vertices {

		dependencies := vertex.dependencies

		for _, dependency := range dependencies {

			v := graph.vertices[dependency]

			for _, dep := range v.dependencies {

				if dep == target {

					fmt.Fprintf(os.Stderr, "cyclic dependency is detected %s\n", target)
					os.Exit(-1)
				}

			}

		}

	}

}

// CheckTarget checks that the input target in the command line exists in the graph representing the makefile
func (graph *Graph) CheckTarget(t string) error {

	for target := range graph.vertices {

		if t == target {
			return nil
		}
	}
	return ErrTargetNotFound
}

// OrderOfExecution it ensures that the commands of the dependencies are executed first before the commands of the input target
func (graph *Graph) OrderOfExecution(t string) {

	target := graph.vertices[t]
	dependencies := target.dependencies

	// Check if there are any dependencies
	if len(dependencies) > 0 {
		for _, dependency := range dependencies {
			dep := graph.vertices[dependency]

			// execute the commands of the dependencies first
			for _, cmd := range dep.cmds {
				ExecCommand(cmd)
			}
		}
	}

	// execute the commands of the target
	for _, cmd := range target.cmds {

		ExecCommand(cmd)

	}

}
