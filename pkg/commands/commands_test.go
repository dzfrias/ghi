package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrNotEnoughArgs(t *testing.T) {
	err := errNotEnoughArgs{"test"}
	want := "not enough args to command `test`"
	assert.Equal(t, err.Error(), want)
}

// func TestCurrentRepo(t *testing.T) {
// 	assert.Equal(t, repo{"dzfrias", "ghi"}, currentRepo())
// }

func TestCurrentRepoNoOrigin(t *testing.T) {
	ex, err := os.Executable()
	if err != nil {
		msg := fmt.Sprintf("Could not get user home directory: %v", err)
		panic(msg)
	}
	os.Chdir(filepath.Dir(ex))
	assert.Equal(t, repo{"", ""}, currentRepo())
}

func TestNewRepo(t *testing.T) {
	r, err := newRepo("dzfrias/ghi")
	assert.Nil(t, err)
	assert.Equal(t, repo{"dzfrias", "ghi"}, r)
}

func TestNewRepoBadRepo(t *testing.T) {
	r, err := newRepo("testing")
	assert.ErrorIs(t, errInvalidRepo, err)
	assert.Equal(t, repo{"", ""}, r)
}

func TestRepoString(t *testing.T) {
	r := repo{"dzfrias", "ghi"}
	assert.Equal(t, "dzfrias/ghi", r.String())
}
