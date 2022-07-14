package issues

import (
	"errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/dzfrias/ghi/pkg/postutil"
)

const ClientId = "3cb4616362f3ae823872"

var ConfigPath string

// Modified during testing
var PollUrl = "https://github.com/login/oauth/access_token"

const pollGrant = "urn:ietf:params:oauth:grant-type:device_code"

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Could not get user home dir")
	}
	// Path to user's credentials
	ConfigPath = path.Join(home, ".config", "ghi", "config")
}

// CredsPoll requests for the user's access token at the given interval (seconds)
func CredsPoll(i int, code string) (url.Values, error) {
	v := url.Values{
		"client_id":   {ClientId},
		"device_code": {code},
		"grant_type":  {pollGrant},
	}
	for {
		// Send POST request and get the url encoded response as map
		respVals, err := postutil.PostParse(PollUrl, v)
		if err != nil {
			return nil, err
		}
		switch respVals.Get("error") {
		case "authorization_pending":
			time.Sleep(time.Second * time.Duration(i))
		case "expired_token":
			return nil, errors.New("access token expired, try again")
		default:
			return respVals, nil
		}
	}
}

// StoreCreds stores the given credentials in ~/.config/ghi/config
func StoreCreds(creds string, loc string) error {
	locDir := filepath.Dir(loc)
	if err := os.MkdirAll(locDir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(loc)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(creds + "\n")
	if err != nil {
		return err
	}

	return nil
}

// GetCreds gets the user's authorized credentials
func GetCreds(fname string) (string, error) {
	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) {
		return "", errors.New("no credentials. Run `ghi auth` to access this feature")
	}

	creds, err := ioutil.ReadFile(fname)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(creds)), err
}
