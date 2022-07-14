package commands

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

const deviceCode = "testing"
const userCode = "1234"
const accessToken = "testing"

func TestLogin(t *testing.T) {
	const credsPath = "./creds.txt"

	loginS := setupLoginServer(t)
	loginUrl = loginS.URL
	pollS := setupPollServer(t)
	issues.PollUrl = pollS.URL
	issues.ConfigPath = credsPath
	defer pollS.Close()
	defer loginS.Close()

	got := testutil.CapStdout(&out, func() {
		args := os.Args[0:1]
		args = append(args, "login")
		err := app.Run(args)
		assert.Nil(t, err)
	})
	want := fmt.Sprintf("Go to https://github.com/login/device and enter this code: %s\n", userCode)
	assert.Equal(t, want, got)
	assert.Equal(t, accessToken+"\n", testutil.Readfile(credsPath))
	os.Remove(credsPath)
}

func TestLoginCredsExist(t *testing.T) {
	const credsPath = "./creds.txt"
	_, err := os.Create(credsPath)
	if err != nil {
		panic(err)
	}
	issues.ConfigPath = credsPath
	got := testutil.CapStdout(&out, func() {
		args := os.Args[0:1]
		args = append(args, "login")
		err := app.Run(args)
		assert.Nil(t, err)
	})
	want := "Credentials already exist.\n"
	assert.Equal(t, want, got)
}

func setupLoginServer(t *testing.T) *httptest.Server {
	resp := url.Values{
		"user_code":   {userCode},
		"interval":    {"0"},
		"device_code": {deviceCode},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := &r.Form
		assert.Equal(t, issues.ClientId, form.Get("client_id"))
		assert.Equal(t, "repo", form.Get("scope"))
		w.Write([]byte(resp.Encode()))
	}))
	return s
}

func setupPollServer(t *testing.T) *httptest.Server {
	resp := url.Values{
		"access_token": {accessToken},
	}
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := &r.Form
		assert.Equal(t, deviceCode, form.Get("device_code"))
		w.Write([]byte(resp.Encode()))
	}))
	return s
}
