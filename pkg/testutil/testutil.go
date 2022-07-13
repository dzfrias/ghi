package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/urfave/cli/v2"
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

// CapStdout captures the stdout of the provided function
func CapStdout(out *io.Writer, f func()) string {
	*out = new(bytes.Buffer)
	f()
	got := (*out).(*bytes.Buffer).String()
	return got
}

// NewApp makes mock cli app (mirroring ghi's interface)
func NewApp(listf func(*cli.Context) error) *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			{
				Name:   "list",
				Action: listf,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "page",
						Value: 1,
					},
				},
			},
		},
	}
}
