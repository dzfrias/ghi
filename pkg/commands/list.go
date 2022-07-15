package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/urfave/cli/v2"
)

// Modified during testing
var out io.Writer = os.Stdout

// List lists the issues with an optional query
func List(ctx *cli.Context) error {
	var result *issues.IssuesSearchResult
	var err error

	page := ctx.Int("page")
	if page < 1 {
		return errors.New("page number must be greater than 0")
	}
	if arg1 := ctx.Args().Get(0); arg1 == "" {
		curRep := "repo:" + currentRepo().String()
		// Default query
		result, err = issues.Search(curRep+" is:open", page)
	} else {
		result, err = issues.Search(
			strings.Join(ctx.Args().Slice(), " "), page)
	}
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Fprintf(out, "#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
	if len(result.Items) < result.TotalCount {
		issNum := (20 * (page - 1)) + 1
		fmt.Fprintf(out, "(Showing issues %d-%d)\n", issNum, issNum+19)
	}

	return nil
}
