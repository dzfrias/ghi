package commands

import (
	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/urfave/cli/v2"
)

// Open opens a issue in GitHub (in cli form)
func Open(ctx *cli.Context) error {
	var r repo

	if repoStr := ctx.Args().Get(0); repoStr == "" {
		r = currentRepo()
	} else {
		var err error
		r, err = newRepo(repoStr)
		if err != nil {
			return err
		}
	}

	issues.Open(r.Owner, r.Name)
	return nil
}
