// Ghi interacts with GitHub Issues
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/dzfrias/ghi/issues"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:      "ghi",
		Usage:     "Interact with GitHub Issues from the command line",
		UsageText: "ghi [global options] <command> [command options] [arguments...]",
		Commands: []*cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"l"},
				Usage:     "Lists issues",
				UsageText: "ghi list [command options] [search terms...]",
				Action:    list,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// list() lists the issues with an optional query
func list(ctx *cli.Context) error {
	var result *issues.IssuesSearchResult
	var err error
	if arg1 := ctx.Args().Get(0); arg1 == "" {
		result, err = issues.SearchIssues(repoQuery() + " is:open")
	} else {
		result, err = issues.SearchIssues(
			strings.Join(ctx.Args().Slice(), " "))
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
		fmt.Println("(Showing issues 1-20)")
	}

	return nil
}

// repoQuery() gets the current repo and puts it into the query format
func repoQuery() string {
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

	ori, err := exec.Command(
		"sh",
		"-c",
		fmt.Sprintf("cd %s; git config --get remote.origin.url", pwd),
	).Output()
	if err != nil {
		return ""
	}

	repo := strings.TrimPrefix(string(ori), "https://github.com/")

	return "repo:" + repo[:len(repo)-5]
}
