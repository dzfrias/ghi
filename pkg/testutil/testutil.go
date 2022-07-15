// Package testutil provides utilities for testing
package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/urfave/cli/v2"
)

// SetupIssueServer sets up a local server that serves fake data, simulating
// https://api.github.com/search/issues
func SetupIssueServer(data any) *httptest.Server {
	// Start a local HTTP server
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		v := r.URL.Query()
		q := v.Get("q")
		// Simulate invalid query
		if q == "invalidRepo" {
			w.WriteHeader(http.StatusUnprocessableEntity)
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
func NewApp(listf, loginf, closef func(*cli.Context) error) *cli.App {
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
			{
				Name:   "login",
				Action: loginf,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "force",
						Value: false,
					},
				},
			},
			{
				Name:   "close",
				Action: closef,
			},
		},
	}
}

// Reafile reads a file, thowing a fatal error if the read fails
func Readfile(fname string) string {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}

	return string(b)
}

// LoadJson loads a json file into a struct
func LoadJson(fname string, res any) {
	b, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(b, &res)
}

// MakeDummyFile creates a file with fake data in it
func MakeDummyFile(name string) {
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString("testing\n")
	if err != nil {
		panic(err)
	}
}
