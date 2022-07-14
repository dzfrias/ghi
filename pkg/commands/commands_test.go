package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotEnoughArgs(t *testing.T) {
	err := errNotEnoughArgs{"test"}
	want := "not enough args to command `test`"
	assert.Equal(t, err.Error(), want)
}
