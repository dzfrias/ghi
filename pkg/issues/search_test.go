package issues

import (
	"os"
	"testing"

	"github.com/dzfrias/ghi/pkg/testutil"
	"github.com/stretchr/testify/assert"
)

var data IssuesSearchResult

func init() {
	testutil.LoadJson("../../testdata/issues.json", &data)
}

func TestSearch(t *testing.T) {
	v, err := Search("repo", 1)
	assert.Nil(t, err)
	assert.Equal(t, v, &data)
}

func TestInvalidSearch(t *testing.T) {
	_, err := Search("invalidRepo", 1)
	const want = "search failed: Invalid search terms or the repository has no issues"
	assert.Equal(t, want, err.Error())
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
