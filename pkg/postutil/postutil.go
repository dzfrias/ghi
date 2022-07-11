// Package postutil provides an easy http POST request abstraction
package postutil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// PostParse sends a POST request and parses the query response as key-value pairs
func PostParse(s string, data url.Values) (url.Values, error) {
	resp, err := http.PostForm(s, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login request failed: %s", resp.Status)
	}

	vals, err := parseQuery(resp.Body)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// parseQuery returns a map of key-value pairs in a url query
func parseQuery(r io.ReadCloser) (url.Values, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	u, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}

	return u, nil
}
