package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var AutumnVersion = "dev"

func version(c *cli.Context) error {
	fmt.Println(AutumnVersion)
	return nil
}
