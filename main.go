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
		Name:  "ghi",
		Usage: "Interact with GitHub Issues from the command line",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List issues",
				Action:  list,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "query",
						Aliases: []string{"q"},
						Usage:   "Search using `TERMS`",
						Value:   currentRepo() + " is:open",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func list(ctx *cli.Context) error {
	result, err := issues.SearchIssues(strings.Split(ctx.String("q"), " "))
	if err != nil {
		return err
	}
	fmt.Printf("%d issues:\n", len(result))
	for _, item := range result {
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number, item.User.Login, item.Title)
	}

	return nil
}

func currentRepo() string {
	pwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	_, err = exec.Command("sh", "-c", fmt.Sprintf("cd %s; git status", pwd)).Output()
	if err != nil {
		return ""
	}

	ori, err := exec.Command("sh", "-c", fmt.Sprintf("cd %s; git config --get remote.origin.url", pwd)).Output()
	if err != nil {
		return ""
	}

	repo := strings.TrimPrefix(string(ori), "https://github.com/")

	return "repo:" + repo[:len(repo)-5]
}
