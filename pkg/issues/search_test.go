package issues

import (
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

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
	code := runTests(m)
	os.Exit(code)
}

func runTests(m *testing.M) int {
	server := testutil.SetupIssueServer(data)
	defer server.Close()
	IssuesURL = server.URL
	return m.Run()
}
