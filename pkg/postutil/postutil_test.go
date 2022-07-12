package postutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var resp = url.Values{
	"testing":  {"test"},
	"testing2": {"test2"},
}

func TestPostParse(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		p := r.Form.Get("testing")
		assert.NotEqual(t, p, "")
		w.Write([]byte(resp.Encode()))
	}))
	defer server.Close()

	v, err := PostParse(server.URL, url.Values{"testing": {"test"}})
	assert.Nil(t, err)
	assert.Equal(t, v, resp)
}
