// Ghi interacts with GitHub Issues from the command line
package main

import (
	"log"
	"os"

	"github.com/dzfrias/ghi/pkg/commands"
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
				ArgsUsage: "[search terms...]",
				Action:    commands.List,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "page",
						Aliases: []string{"p"},
						Value:   1,
						Usage:   "Page of results",
					},
				},
			},
			{
				Name:      "login",
				Usage:     "Brings up login to have advanced options",
				UsageText: "ghi login [command options]",
				Action:    commands.Login,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "force",
						Aliases: []string{"f"},
						Value:   false,
						Usage:   "Create credentials even if they already exist",
					},
				},
			},
			{
				Name:      "close",
				Usage:     "Closes an issue (needs login)",
				UsageText: "ghi close [command options] <issue> <repo>",
				Action:    commands.Close,
			},
			{
				Name:      "open",
				Usage:     "Opens an issue (needs login)",
				UsageText: "ghi open [command options] <repo>",
				Action:    commands.Open,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
