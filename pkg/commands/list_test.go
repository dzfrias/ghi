package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var data issues.IssuesSearchResult

func init() {
	testutil.LoadJson("../../testdata/issues.json", &data)
}

var app = testutil.NewApp(
	List,
	Login,
	Close,
)

func TestList(t *testing.T) {
	const expected = `21 issues:
#1      TestUser Testing
#2      TestUser Testing
#3      TestUser Testing
#4      TestUser Testing
#5      TestUser Testing
#6      TestUser Testing
#7      TestUser Testing
#8      TestUser Testing
#9      TestUser Testing
#10     TestUser Testing
#11     TestUser Testing
#12     TestUser Testing
#13     TestUser Testing
#14     TestUser Testing
#15     TestUser Testing
#16     TestUser Testing
#17     TestUser Testing
#18     TestUser Testing
#19     TestUser Testing
#20     TestUser Testing
(Showing issues 1-20)
`
	for _, arg := range []string{"repo", ""} {
		got := testutil.CapStdout(&out, func() {
			args := os.Args[0:1]
			args = append(args, "list", arg)
			err := app.Run(args)
			assert.Nil(t, err)
		})
		assert.Equal(t, expected, got)
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
		msg := fmt.Sprintf("Could not get user home directory: %v", err)
		panic(msg)
	}
	os.Chdir(filepath.Dir(ex))
	assert.Equal(t, "", currentRepo())
}

func TestMain(m *testing.M) {
	code := runTests(m)
	os.Exit(code)
}

func runTests(m *testing.M) int {
	server := testutil.SetupIssueServer(data)
	defer server.Close()
	issues.IssuesURL = server.URL
	return m.Run()
}
