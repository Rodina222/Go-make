package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMakeFile(t *testing.T) {

	graph := NewGraph()

	t.Parallel()

	expectedGraph := Graph{
		vertices: map[string]Vertex{
			"build": {
				dependencies: []string{},
				cmds:         []string{"@echo 'executing build'", "echo 'cmd2'"},
			},
			"test": {
				dependencies: []string{"build"},
				cmds:         []string{"@echo 'executing test'"},
			},
			"publish": {
				dependencies: []string{"test"},
				cmds:         []string{"@echo 'executing publish'"},
			},
		},
	}

	t.Run("valid makefile", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)

		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)
		assert.NoError(t, err)

		if !reflect.DeepEqual(graph, expectedGraph) {
			t.Errorf("graph %v does not match expected graph %v", graph, expectedGraph)
		}

	})

	t.Run("invalid makefile", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile2")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)

		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)

		assert.Equal(t, ErrInvalidFormat, err, "want error %q but got %q", ErrInvalidFormat, err)
	})

}

func TestCheckCmds(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	t.Run("all targets have commands", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)
		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)
		assert.NoError(t, err)

		err = graph.CheckCmds()
		assert.NoError(t, err)
	})

	t.Run("a target or more have no commands", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile3")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)
		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)
		assert.NoError(t, err)

		err = graph.CheckCmds()
		assert.Error(t, err)
	})

}

func TestCyclicDependency(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	t.Run("no cyclic dependency exists", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)
		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)
		assert.NoError(t, err)

		err = graph.CheckCyclicDependency()
		assert.NoError(t, err)
	})

	t.Run("cyclic dependency exists", func(t *testing.T) {

		dir := t.TempDir()

		filePath := filepath.Join(dir, "Makefile4")

		err := os.WriteFile(filePath, make([]byte, 0), 0644)
		assert.NoError(t, err)

		filename := filepath.Base(filePath)

		err = graph.ParseMakeFile(filename)
		assert.NoError(t, err)

		err = graph.CheckCyclicDependency()
		assert.Equal(t, ErrCyclicDependency, err, "want error %q but got %q", ErrCyclicDependency, err)
	})

}

func TestOrderOfExecution(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	dir := t.TempDir()

	filePath := filepath.Join(dir, "Makefile")

	err := os.WriteFile(filePath, make([]byte, 0), 0644)
	assert.NoError(t, err)

	filename := filepath.Base(filePath)

	err = graph.ParseMakeFile(filename)
	assert.NoError(t, err)

	vertex := graph.vertices["build"]
	want := vertex.cmds

	got := graph.OrderOfExecution("build")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("order of cmds of got %s does not match want %v", got, want)
	}

}
