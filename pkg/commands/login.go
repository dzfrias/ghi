package commands

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/dzfrias/ghi/pkg/postutil"
	"github.com/urfave/cli/v2"
)

// Login puts the user through the authentication process to store credentials
func Login(ctx *cli.Context) error {
	if _, err := issues.GetCreds(); err == nil && !ctx.Bool("force") {
		fmt.Println("Credentials already exist.")
		return nil
	}
	const loginUrl = "https://github.com/login/device/code"

	data := url.Values{
		"client_id": {issues.ClientId},
		"scope":     {"repo"},
	}

	// Request a login and get url encoded response
	vals, err := postutil.PostParse(loginUrl, data)
	if err != nil {
		return err
	}
	fmt.Printf("Go to https://github.com/login/device and enter this code: %s\n", vals.Get("user_code"))

	i, err := strconv.Atoi(vals.Get("interval"))
	if err != nil {
		return err
	}
	code := vals.Get("device_code")

	// Try to get auth info every i seconds
	authInfo, err := issues.CredsPoll(i, code)
	if err != nil {
		return err
	}

	issues.StoreCreds(authInfo.Get("access_token"))

	return nil
}
