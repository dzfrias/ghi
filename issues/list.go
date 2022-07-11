// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Modifications to the original have been made
// Source: https://github.com/adonovan/gopl.io/blob/master/ch4/github/search.go

package issues

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(term string) ([]*Issue, error) {
	q := url.QueryEscape(term)

	req, err := http.NewRequest("GET", IssuesURL+"?q="+q, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(
		"Accept", "application/vnd.github.v3.text-match+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == 422 {
		return nil, fmt.Errorf("search failed: Invalid search terms or the repository has no issues")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result struct{ Items []*Issue }
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.Items, nil
}
