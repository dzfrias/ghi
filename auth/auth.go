// Package auth handles authentication for the user
package auth

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

const clientId = "3cb4616362f3ae823872"

var configPath string

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln("Could not get user home dir")
	}
	configPath = path.Join(home, ".config", "ghi", "config")
}

// Auth puts the user through the authentication process to store credentials
func Auth(ctx *cli.Context) error {
	if credsExist() {
		fmt.Println("Credentials already exist.")
		return nil
	}
	const loginUrl = "https://github.com/login/device/code"

	data := url.Values{
		"client_id": {clientId},
		"scope":     {"repo"},
	}

	vals, err := postParse(loginUrl, data)
	if err != nil {
		return err
	}
	fmt.Printf("Go to https://github.com/login/device and enter this code: %s\n", vals.Get("user_code"))

	i, err := strconv.Atoi(vals.Get("interval"))
	if err != nil {
		return err
	}
	code := vals.Get("device_code")

	authInfo, err := poll(i, code)
	if err != nil {
		return err
	}

	store(authInfo.Get("access_token"))

	return nil
}

func postParse(s string, data url.Values) (url.Values, error) {
	resp, err := http.PostForm(s, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login request failed: %s", resp.Status)
	}

	vals, err := parseUrl(resp.Body)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

func parseUrl(r io.ReadCloser) (url.Values, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	u, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func poll(i int, code string) (url.Values, error) {
	const pollUrl = "https://github.com/login/oauth/access_token"
	const grantType = "urn:ietf:params:oauth:grant-type:device_code"
	v := url.Values{
		"client_id":   {clientId},
		"device_code": {code},
		"grant_type":  {grantType},
	}
	for {
		respVals, err := postParse(pollUrl, v)
		if err != nil {
			return nil, err
		}
		if respVals.Get("error") != "authorization_pending" {
			return respVals, nil
		}
		time.Sleep(time.Second * time.Duration(i))
	}
}

func store(creds string) error {
	conDir := filepath.Dir(configPath)
	if err := os.MkdirAll(conDir, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(configPath)
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
func GetCreds() (string, error) {
	if !credsExist() {
		return "", errors.New("no credentials. Run `ghi auth` to access this feature")
	}

	creds, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(creds)), err
}

func credsExist() bool {
	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}
