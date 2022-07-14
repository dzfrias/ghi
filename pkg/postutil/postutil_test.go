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
		assert.NotEmpty(t, p)
		w.Write([]byte(resp.Encode()))
	}))
	defer server.Close()

	v, err := PostParse(server.URL, url.Values{"testing": {"test"}})
	assert.Nil(t, err)
	assert.Equal(t, v, resp)
}

func TestPostParseFail(t *testing.T) {
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(""))
	}))
	defer server.Close()

	const want = "POST request failed: 500 Internal Server Error"
	_, err := PostParse(server.URL, url.Values{"none": {"none"}})
	assert.Equal(t, want, err.Error())
}

func TestEncodeMap(t *testing.T) {
	m := map[string]string{"testing": "test"}
	b, err := EncodeMap(m)
	assert.Nil(t, err)
	assert.Equal(t, `{"testing":"test"}`+"\n", b.String())
}
