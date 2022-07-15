package issues

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	const user = "dzfrias"
	const repo = "ghi"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		acc := r.Header.Get("Accept")
		assert.Equal(t, "application/vnd.github+json", acc)
		assert.Equal(t, "token testing", auth)
		u := strings.Split(r.URL.String(), "/")[1:]
		assert.Equal(t, user, u[0])
		assert.Equal(t, repo, u[1])
		bod, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)

		const want = `{"body":"test","title":"testing"}` + "\n"
		assert.Equal(t, want, string(bod))
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	OpenUrl = server.URL
	ConfigPath = "./creds_testdata.txt"
	testutil.MakeDummyFile(ConfigPath)
	defer os.Remove(ConfigPath)

	const issFile = "./issue_test.txt"
	makeDummyIssF(issFile, "test")

	err := openFrom(issFile, user, repo)
	assert.Nil(t, err)
}

func TestOpenReq(t *testing.T) {
	ConfigPath = "./creds_testdata.txt"
	testutil.MakeDummyFile(ConfigPath)
	defer os.Remove(ConfigPath)
	var iss = issueFile{"testing", "test"}

	req, err := openReq("localhost:8000", iss)
	assert.Nil(t, err)
	head := req.Header
	assert.Equal(t, "token testing", head.Get("Authorization"))
	assert.Equal(t, "application/vnd.github+json", head.Get("Accept"))

	assert.Equal(t, "localhost:8000", req.URL.String())
}

func TestReadLines(t *testing.T) {
	const fname = "./test.txt"
	makeDummyIssF(fname, "testing\ntest")
	defer os.Remove(fname)
	l, err := readLines(fname)
	assert.Nil(t, err)

	want := []string{"testing", "testing", "test"}
	assert.Equal(t, want, l)
}

func TestReadIssF(t *testing.T) {
	const issFile = "./issue_test.txt"
	makeDummyIssF(issFile, "test")
	defer os.Remove(issFile)

	iss, err := readIssF(issFile)
	assert.Nil(t, err)
	assert.Equal(t, iss, issueFile{"testing", "test"})
}

func TestReadIssFNoCont(t *testing.T) {
	const issFile = "./issue_test.txt"
	f, err := os.Create(issFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	defer os.Remove(issFile)

	_, err = readIssF(issFile)
	assert.ErrorIs(t, ErrBadTempFile, err)
}

func makeDummyIssF(fname, cont string) {
	testutil.MakeDummyFile(fname)
	f, err := os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(cont); err != nil {
		panic(err)
	}
}
