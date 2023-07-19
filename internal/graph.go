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
	ErrTargetNotFound = errors.New("target is not found")

	// ErrInvalidFormat is returned when the format of the input makefile is invalid
	ErrInvalidFormat = errors.New("format is invalid")

	// ErrNoCommandsFound is returned when there is a target that has no commands to be executed
	ErrNoCommandsFound = errors.New("no commands to be executed for this target")

	// ErrCyclicDependency is returned once a cyclic dependency is detected in the graph
	ErrCyclicDependency = errors.New("cyclic dependency is detected")
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
		return err
	}

	defer content.Close()

	// create a scanner to read the file line by line
	scanner := bufio.NewScanner(content)

	if err := scanner.Err(); err != nil {
		return errors.New("scanner failed to scan the file")
	}

	// create a vertex in the graph
	vertex := NewVertex()

	// scanning the file line by line
	for scanner.Scan() {

		line := scanner.Text()

		// skipping empty lines and comments
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "\t") && target != "" {

			cmd := strings.TrimPrefix(line, "\t")
			vertex.cmds = append(vertex.cmds, cmd)
			graph.vertices[target] = vertex

			continue

		}

		if strings.Count(line, ":") == 1 {

			// create a new vertex in the graph
			vertex = NewVertex()

			// specify the position of colon  ":"
			colonIndex := strings.Index(line, ":")

			//extract the target from the line
			target = strings.TrimSpace(line[:colonIndex])

			if target == "" {
				return ErrInvalidFormat
			}

			// add the vertex to the graph
			graph.vertices[target] = vertex

			// add dependencies (if found)
			if len(line) > colonIndex+1 {

				deps := strings.TrimSpace(line[colonIndex+1:])
				vertex.AddDependencies(deps)

			}
			continue

		}
		return ErrInvalidFormat
	}
	return nil
}

// Execute is a method that checks all the other methods on the graph and returns an error once one of them returns an error
func (graph *Graph) Execute(target string) error {

	// check cyclic dependency
	err := graph.CheckCyclicDependency()

	if err != nil {
		fmt.Println(err)
		return err
	}

	//check all targets have commands
	err = graph.CheckCmds()

	if err != nil {
		fmt.Println(err)
		return err
	}

	// get the order of executing commands and execute them
	err = graph.ExecuteInOrder(target)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// CheckCmds loops on all targets and checks that all of them have commands to be executed, otherwise it will return error
func (graph *Graph) CheckCmds() error {

	for target, vertex := range graph.vertices {

		if len(vertex.cmds) == 0 {

			return fmt.Errorf("%w : %s", ErrNoCommandsFound, target)

		}
	}

	return nil
}

// CheckCyclicDependency checks if a cyclic dependency exists or not and returns error once it is found
func (graph *Graph) CheckCyclicDependency() error {

	visited := make(map[string]bool)
	onStack := make(map[string]bool)

	for target := range graph.vertices {
		if !visited[target] {

			// perform DFS fo every unvisited target
			if graph.hasCycle(target, visited, onStack) {
				return ErrCyclicDependency
			}
		}
	}

	return nil
}

// hasCycle is a helper function for CheckCyclicDependency that performs DFS
func (graph *Graph) hasCycle(v string, visited map[string]bool, onStack map[string]bool) bool {

	visited[v] = true
	onStack[v] = true

	for _, dependency := range graph.vertices[v].dependencies {

		if !visited[dependency] {
			return graph.hasCycle(dependency, visited, onStack)
		}

		if onStack[dependency] {
			return true
		}
	}

	onStack[v] = false
	return false
}

// ExecuteInOrder ensures that the commands of the dependencies are executed first before the commands of the input target
func (graph *Graph) ExecuteInOrder(t string) error {

	target := graph.vertices[t]
	dependencies := target.dependencies

	// Check if there are any dependencies

	if len(dependencies) > 0 {
		for _, dependency := range dependencies {

			// apply recursion to cover all layers of dependencies
			err := graph.ExecuteInOrder(dependency)
			if err != nil {
				return err
			}
		}
	}

	// execute the commands of the target
	for _, cmd := range target.cmds {
		err := ExecCommand(cmd)
		if err != nil {
			return err
		}
	}

	return nil

}
