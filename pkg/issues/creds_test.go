package issues

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

const credsFile = "./testcreds.txt"

func TestGetCreds(t *testing.T) {
	newF := "./creds_test.txt"
	testutil.MakeCredsFile(newF)
	defer os.Remove(newF)
	creds, err := GetCreds(newF)
	assert.Equal(t, creds, "testing", "does not read creds correctly")
	assert.Nil(t, err, "throws an error")
}

func TestGetCredsNoneExist(t *testing.T) {
	_, err := GetCreds(credsFile)
	assert.ErrorIs(t, err, ErrNoCreds)
}

func TestStoreCredsNoConfigDir(t *testing.T) {
	const newCreds = "./.config/ghi/testcreds.txt"
	err := StoreCreds("testing", newCreds)
	assert.Nil(t, err, "StoreCreds with no config dir throws error")
	assert.Equal(t, testutil.Readfile(newCreds), "testing\n")
	os.Remove(newCreds)
}

func TestStoreCreds(t *testing.T) {
	err := StoreCreds("testing", credsFile)
	assert.Nil(t, err, "StoreCreds throws error")
	assert.Equal(t, testutil.Readfile(credsFile), "testing\n")
	os.Remove(credsFile)
}

func TestCredsPoll(t *testing.T) {
	var errResp = url.Values{
		"error": {"authorization_pending"},
	}
	var resp = url.Values{
		"testing": {"success"},
	}
	const wantDc = "1234"
	reqs := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		form := &r.Form
		id := form.Get("client_id")
		dc := form.Get("device_code")
		gt := form.Get("grant_type")
		assert.Equal(t, ClientId, id)
		assert.Equal(t, wantDc, dc)
		assert.Equal(t, pollGrant, gt)
		if reqs == 0 {
			// Simulate waiting for user input
			w.Write([]byte(errResp.Encode()))
			reqs++
		} else {
			w.Write([]byte(resp.Encode()))
		}
	}))
	defer server.Close()

	PollUrl = server.URL
	v, err := CredsPoll(0, wantDc)
	assert.Nil(t, err)
	assert.Equal(t, resp, v)
}

func TestCredsPollExpiredToken(t *testing.T) {
	const msg = "access token expired, try again"
	var resp = url.Values{
		"error": {"expired_token"},
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(resp.Encode()))
	}))
	defer server.Close()

	PollUrl = server.URL
	_, err := CredsPoll(0, "none")
	assert.Equal(t, msg, err.Error())
}
