package commands

import (
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var data = issues.IssuesSearchResult{
	TotalCount: 1,
	Items: []*issues.Issue{
		{
			Number:  1,
			HTMLURL: "none",
			Title:   "Testing",
			State:   "open",
			User: &issues.User{
				Login:   "TestUser",
				HTMLURL: "none",
			},
		},
	},
}

var app = testutil.NewApp(List)

func TestList(t *testing.T) {
	got := testutil.CapStdout(&out, func() {
		args := os.Args[0:1]
		args = append(args, "list", "repo")
		err := app.Run(args)
		assert.Nil(t, err)
	})
	target := "1 issues:\n#1      TestUser Testing\n"
	assert.Equal(t, target, got)
}

func TestListBadPage(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "list", "--page", "0", "repo")
	err := app.Run(args)
	assert.Error(t, err)
}

func TestMain(m *testing.M) {
	server := testutil.SetupIssueServer(data)
	defer server.Close()
	issues.IssuesURL = server.URL
	os.Exit(m.Run())
}
