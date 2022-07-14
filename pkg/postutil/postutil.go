// Package postutil provides easy http POST request abstractions
package postutil

import (
	"bytes"
	"encoding/json"
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
		return nil, fmt.Errorf("POST request failed: %s", resp.Status)
	}

	// Parse url encoded key-value pairs
	vals, err := parseQuery(resp.Body)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// EncodeMap encodes a map into a json bytes buffer
func EncodeMap(m map[string]string) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(m)
	if err != nil {
		return buf, err
	}
	return buf, nil
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
