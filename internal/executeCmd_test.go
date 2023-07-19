package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecCommand(t *testing.T) {

	t.Run("valid command", func(t *testing.T) {

		err := ExecCommand("@echo 'executing publish'")
		assert.NoError(t, err)
	})

	t.Run("invalid command", func(t *testing.T) {

		err := ExecCommand("@echoo 'executing publish'")
		assert.Error(t, err)
	})

}
