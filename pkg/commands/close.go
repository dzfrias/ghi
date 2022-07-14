package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/urfave/cli/v2"
)

var errInvalidRepo = errors.New("invalid repo name, must have a '/'")

// Close closes an issue
func Close(ctx *cli.Context) error {
	args := ctx.Args()

	issNum := args.Get(0)
	if issNum == "" {
		return errNotEnoughArgs{"close"}
	}
	fullRepo := args.Get(1)
	if fullRepo == "" {
		fullRepo = currentRepo()[5:]
	}
	fmt.Println(fullRepo)

	split := strings.Split(fullRepo, "/")
	if len(split) < 2 {
		return errInvalidRepo
	}
	o := split[0]
	r := split[1]

	err := issues.CloseIssue(issNum, o, r)
	if err != nil {
		return err
	}

	return nil
}
