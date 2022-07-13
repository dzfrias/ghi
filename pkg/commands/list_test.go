package commands

import (
	"log"
	"os"
	"path/filepath"
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
	for _, arg := range []string{"repo", ""} {
		got := testutil.CapStdout(&out, func() {
			args := os.Args[0:1]
			args = append(args, "list", arg)
			err := app.Run(args)
			assert.Nil(t, err)
		})
		target := "1 issues:\n#1      TestUser Testing\n"
		assert.Equal(t, target, got)
	}
}

func TestListBadPage(t *testing.T) {
	args := os.Args[0:1]
	args = append(args, "list", "--page", "0", "repo")
	err := app.Run(args)
	assert.Error(t, err)
}

func TestCurrentRepo(t *testing.T) {
	assert.Equal(t, "repo:dzfrias/ghi", currentRepo())
}

func TestCurrentRepoNoOrigin(t *testing.T) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Could not get user home directory: %v", err)
	}
	os.Chdir(filepath.Dir(ex))
	assert.Equal(t, "", currentRepo())
}

func TestMain(m *testing.M) {
	server := testutil.SetupIssueServer(data)
	defer server.Close()
	issues.IssuesURL = server.URL
	os.Exit(m.Run())
}
