package main

// DO NOT MODIFY: generated by github.com/ttacon/toml2cli
// If you need to make changes, update the toml config and then regenerate
// this file.

import (
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "autumn"
	app.Usage = "A code generation tool that helps you relax"

	app.Commands = []*cli.Command{
		&cli.Command{
			Description: "Print out the current version",
			Name:        "version",
			Action:      version,
		},
		&cli.Command{
			Name:        "init",
			Description: "Initialize an autumn configuration file",
			Action:      initCommand,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name: "force",
					Aliases: []string{
						"f",
					},
				},
			},
		},
		&cli.Command{
			Action:      get,
			Description: "Retrieve all frameworks.",
			Name:        "get",
		},
		&cli.Command{
			Name:        "plan",
			Description: "Plan the code to generate.",
			Action:      plan,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "out",
					Value: "autumn-plan.json",
					Aliases: []string{
						"o",
					},
				},
				&cli.BoolFlag{
					Name: "force",
					Aliases: []string{
						"f",
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
