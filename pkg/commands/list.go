package commands

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dzfrias/ghi/pkg/issues"
	"github.com/urfave/cli/v2"
)

// List lists the issues with an optional query
func List(ctx *cli.Context) error {
	var result *issues.IssuesSearchResult
	var err error

	page := ctx.Int("page")
	if page < 1 {
		return errors.New("page number must be greater than 0")
	}
	if arg1 := ctx.Args().Get(0); arg1 == "" {
		// Default query
		result, err = issues.Search(currentRepo()+" is:open", page)
	} else {
		result, err = issues.Search(
			strings.Join(ctx.Args().Slice(), " "), page)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%d issues:\n", result.TotalCount)
	for _, item := range result.Items {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}
	if len(result.Items) < result.TotalCount {
		issNum := (20 * (page - 1)) + 1
		fmt.Printf("(Showing issues %d-%d)\n", issNum, issNum+19)
	}

	return nil
}

// currentRepo gets the current repo and puts it into the query format
func currentRepo() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	_, err = exec.Command(
		"sh", "-c", fmt.Sprintf("cd %s; git status", pwd),
	).Output()
	if err != nil {
		return ""
	}

	// Get origin url in cwd
	ori, err := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("cd %s; git config --get remote.origin.url", pwd),
	).Output()
	if err != nil {
		return ""
	}

	repo := strings.TrimPrefix(string(ori), "https://github.com/")

	// Strip the '.git' at the end
	return "repo:" + repo[:len(repo)-5]
}
