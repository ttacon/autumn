package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ttacon/autumn/lib/config"
	"github.com/urfave/cli/v2"
)

func plan(c *cli.Context) error {
	// Plan steps:
	//
	//  1. Check if plan with known name exists.
	//    a. If it does, and force flag isn't set, exit.
	//  2. Check if we're using targeting mode.
	//  3. Load config.
	//  4. See if all frameworks have been pulled down.
	//  5. Generate plan (Model -> <Controller, Router, Service> -> Framework -> Templates).
	//  6. Write file to disk.

	outputFileName := c.String("out")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	root := os.DirFS(cwd)

	var conf config.Config
	if data, err := ioutil.ReadFile(
		filepath.Join(
			autumnDir,
			configFileName,
		),
	); err != nil {
		fmt.Println("failed to load autumn config: ", err)
		return err
	} else if _, err := toml.NewDecoder(bytes.NewBuffer(data)).Decode(&conf); err != nil {
		fmt.Println("autumn config is malformed: ", err)
		return err
	}

	return retrieveSourcesForEngine(conf)
}
