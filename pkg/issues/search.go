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
	"strconv"
)

// Search queries the GitHub issue tracker.
func Search(term string, page int) (*IssuesSearchResult, error) {
	q := url.QueryEscape(term)

	reqUrl := IssuesURL + "?q=" + q + "&per_page=20&page=" + strconv.Itoa(page)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github.v3.text-match+json")
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

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
