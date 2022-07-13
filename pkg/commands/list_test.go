package commands

import (
	"bytes"
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
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

var app = &cli.App{
	Commands: []*cli.Command{
		{
			Name:   "list",
			Action: List,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "page",
					Value: 1,
				},
			},
		},
	},
}

func TestList(t *testing.T) {
	out = new(bytes.Buffer)
	args := os.Args[0:1]
	args = append(args, "list", "repo")
	err := app.Run(args)
	got := out.(*bytes.Buffer).String()
	target := "1 issues:\n#1      TestUser Testing\n"
	assert.Equal(t, target, got)
	assert.Nil(t, err)
}

func TestMain(m *testing.M) {
	server := testutil.SetupIssueServer(data)
	defer server.Close()
	issues.IssuesURL = server.URL
	os.Exit(m.Run())
}
