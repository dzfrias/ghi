package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// SetupIssueServer sets up a local server that serves fake data, simulating
// https://api.github.com/search/issues
func SetupIssueServer(data interface{}) *httptest.Server {
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
