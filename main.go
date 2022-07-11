// Ghi interacts with GitHub Issues
package main

import (
	"log"
	"os"

	"github.com/dzfrias/ghi/auth"
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
				Action:    list,
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
				Action:    auth.Auth,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
