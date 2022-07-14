package issues

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dzfrias/ghi/pkg/postutil"
)

var closeUrl = "https://api.github.com/repos"

// CloseIssue closes a GitHub issue
func CloseIssue(id, owner, repo string) error {
	buf, err := postutil.EncodeMap(map[string]string{"state": "closed"})
	if err != nil {
		return err
	}

	reqUrl := fmt.Sprintf("%s/%s/%s/issues/%s", closeUrl, owner, repo, id)
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("close issue failed with status: %s", resp.Status)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respMap := make(map[string]any)
	if err = json.Unmarshal([]byte(string(b)), &respMap); err != nil {
		return err
	}
	if respMap["message"] != nil {
		return fmt.Errorf("close issue failed: %s", respMap["message"])
	}

	return nil
}
