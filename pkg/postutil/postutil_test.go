package postutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var resp = url.Values{
	"testing":  {"test"},
	"testing2": {"test2"},
}

func TestPostParse(t *testing.T) {
	const test = `PostParse(server.URL, url.Values{"testing": {"test"})`
	// Start a local HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		p := r.Form.Get("testing")
		if p == "" {
			t.Errorf(`%s - %q != "test"'`, test, p)
		}
		w.Write([]byte(resp.Encode()))
	}))
	defer server.Close()

	v, err := PostParse(server.URL, url.Values{"testing": {"test"}})
	if err != nil {
		t.Errorf(`%s threw error %v`, test, err)
	}
	if !reflect.DeepEqual(v, resp) {
		t.Errorf(`%s != %v`, test, resp)
	}
}
