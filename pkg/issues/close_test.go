package issues

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	const issNum = "1"
	const user = "dzfrias"
	const repo = "testing"

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

	CloseUrl = server.URL
	ConfigPath = "./creds_testdata.txt"
	err := CloseIssue(issNum, user, repo)
	assert.Nil(t, err)
}

func TestCloseNoCreds(t *testing.T) {
	ConfigPath = "./doesnotexist.txt"
	err := CloseIssue("1", "2", "3")
	assert.ErrorIs(t, err, ErrNoCreds)
}
