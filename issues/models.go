// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Modifications to the original have been made
// Source: https://github.com/adonovan/gopl.io/blob/master/ch4/github/github.go

// Package issues is a Go API for GitHub issues
package issues

import (
	"time"
)

const IssuesURL = "https://api.github.com/search/issues"

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}
