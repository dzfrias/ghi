package issues

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

var data = IssuesSearchResult{
	TotalCount: 1,
	Items: []*Issue{
		{
			Number:  1,
			HTMLURL: "none",
			Title:   "Testing",
			State:   "open",
			User: &User{
				Login:   "TestUser",
				HTMLURL: "none",
			},
		},
	},
}

func TestSearch(t *testing.T) {
	v, err := Search("repo", 1)
	assert.Nil(t, err)
	assert.Equal(t, v, &data)
}

func TestInvalidSearch(t *testing.T) {
	_, err := Search("invalidRepo", 1)
	assert.Error(t, err)
}

func TestMain(m *testing.M) {
	server = setupServer()
	defer server.Close()
	IssuesURL = server.URL
	os.Exit(m.Run())
}

func setupServer() *httptest.Server {
	// Start a local HTTP server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		v := r.URL.Query()
		q := v.Get("q")
		// Simulate invalid query
		if q == "invalidRepo" {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(""))
			return
		}
		json.NewEncoder(w).Encode(data)
	}))
	return s
}
