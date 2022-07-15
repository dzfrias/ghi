package commands

import (
	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/urfave/cli/v2"
)

// Close closes an issue
func Close(ctx *cli.Context) error {
	args := ctx.Args()

	issNum := args.Get(0)
	if issNum == "" {
		return errNotEnoughArgs{"close"}
	}
	var r repo
	if fullRepo := args.Get(1); fullRepo == "" {
		r = currentRepo()
	} else {
		var err error
		r, err = newRepo(fullRepo)
		if err != nil {
			return err
		}
	}

	err := issues.CloseIssue(issNum, r.Owner, r.Name)
	if err != nil {
		return err
	}

	return nil
}
