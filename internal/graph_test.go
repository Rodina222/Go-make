package internal

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)
const validMakefile = `build:
	@echo 'executing build'
	echo 'cmd2'

	test: build
	@echo 'executing test'

	publish: test 
	@echo 'executing publish'`

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

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(validMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.NoError(t, err)

		if !reflect.DeepEqual(graph, expectedGraph) {
			t.Errorf("graph %v does not match expected graph %v", graph, expectedGraph)
		}

	})

	t.Run("invalid makefile", func(t *testing.T) {

		invalidMakefile := `build:
	@echo 'executing build'
	echo 'cmd2'

	  : build
	@echo 'executing test'

	publish: test 
	@echo 'executing publish'`

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(invalidMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.Error(t, err)

		assert.Equal(t, ErrInvalidFormat, err, "want error %q but got %q", ErrInvalidFormat, err)

	})

}

func TestCheckCmds(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	t.Run("all targets have commands", func(t *testing.T) {

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(validMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.NoError(t, err)

		err = graph.CheckCmds()
		assert.NoError(t, err)
	})

	t.Run("a target or more have no commands", func(t *testing.T) {

		invalidMakefile := `build:
		@echo 'executing build'
		echo 'cmd2'
	
		test: build
	
		publish: test 
		@echo 'executing publish'`

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(invalidMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.NoError(t, err)

		err = graph.CheckCmds()
		assert.Error(t, err)

	})

}

func TestCyclicDependency(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	t.Run("no cyclic dependency exists", func(t *testing.T) {

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(validMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.NoError(t, err)

		err = graph.CheckCyclicDependency()
		assert.NoError(t, err)

	})

	t.Run("cyclic dependency exists", func(t *testing.T) {

		invalidMakefile := `build:
		@echo 'executing build'
		echo 'cmd2'
	
		test: build publish
		@echo 'executing test'
	
		publish: test 
		@echo 'executing publish'`

		dir := os.TempDir()
		filePath := filepath.Join(dir, "Makefile")

		file, err := os.Create(filePath)
		assert.NoError(t, err)

		defer os.Remove(file.Name())

		_, err = file.WriteString(invalidMakefile)
		assert.NoError(t, err)

		err = graph.ParseMakeFile(file.Name())
		assert.NoError(t, err)

		err = graph.CheckCyclicDependency()
		assert.Equal(t, ErrCyclicDependency, err, "want error %q but got %q", ErrCyclicDependency, err)

	})

}

func TestOrderOfExecution(t *testing.T) {

	graph := NewGraph()
	t.Parallel()

	dir := os.TempDir()
	filePath := filepath.Join(dir, "Makefile")

	file, err := os.Create(filePath)
	assert.NoError(t, err)

	defer os.Remove(file.Name())

	_, err = file.WriteString(validMakefile)
	assert.NoError(t, err)

	err = graph.ParseMakeFile(file.Name())
	assert.NoError(t, err)

	vertex := graph.vertices["build"]
	want := vertex.cmds

	got := graph.OrderOfExecution("build")

	if !reflect.DeepEqual(got, want) {
		t.Errorf("order of cmds of got %s does not match want %v", got, want)
	}

}
