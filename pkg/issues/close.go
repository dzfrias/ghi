package issues

import (
	"fmt"
	"net/http"

	"github.com/dzfrias/ghi/pkg/postutil"
)

// Modified during testing
var CloseUrl = "https://api.github.com/repos"

// CloseIssue closes a GitHub issue
func CloseIssue(id, owner, repo string) error {
	buf, err := postutil.EncodeMap(map[string]string{"state": "closed"})
	if err != nil {
		panic("error encoding map for closing issue")
	}

	reqUrl := fmt.Sprintf("%s/%s/%s/issues/%s", CloseUrl, owner, repo, id)
	req, err := http.NewRequest(http.MethodPatch, reqUrl, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	s, err := GetCreds(ConfigPath)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "token "+s)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("close issue failed with status: %s", resp.Status)
	}

	return nil
}
