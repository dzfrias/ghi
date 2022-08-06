package commands

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	const issNum = "1"
	const user = "dzfrias"
	const repo = "ghi"

	fullrepo := strings.Join([]string{user, repo}, "/")
	for _, arg := range []string{fullrepo} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			assert.Equal(t, "token testing", auth)
			u := strings.Split(r.URL.String(), "/")[1:]
			assert.Equal(t, user, u[0])
			assert.Equal(t, repo, u[1])
			assert.Equal(t, issNum, u[3])
			bod, err := ioutil.ReadAll(r.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"state":"closed"}`+"\n", string(bod))
		}))
		defer server.Close()

		const credsFile = "test_creds.txt"
		testutil.MakeDummyFile(credsFile)
		defer os.Remove(credsFile)
		issues.CloseUrl = server.URL
		issues.ConfigPath = credsFile

		args := os.Args[0:1]
		args = append(args, "close", issNum, arg)
		err := app.Run(args)
		assert.Nil(t, err)
	}
}

func TestCloseNoCreds(t *testing.T) {
	issues.ConfigPath = "doesnotexist.txt"

	args := os.Args[0:1]
	args = append(args, "close", "1", "dzfrias/ghi")
	err := app.Run(args)
	assert.ErrorIs(t, err, issues.ErrNoCreds)
}

func TestCloseNotEnoughArgs(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "close")
	err := app.Run(args)
	assert.ErrorIs(t, err, errNotEnoughArgs{"close"})
}

func TestCloseInvalidRepo(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "close", "1", "testing")
	err := app.Run(args)
	assert.ErrorIs(t, err, errInvalidRepo)
}
