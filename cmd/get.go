package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ttacon/autumn/lib/config"
	"github.com/ttacon/autumn/lib/engine/retriever"

	"github.com/urfave/cli/v2"
)

func get(c *cli.Context) error {
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

func retrieveSourcesForEngine(c config.Config) error {
	var frameworkGetters = []config.FrameworkGetter{
		c.Controller,
		c.Router,
		c.Service,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	frameworkRetriever, err := retriever.NewFrameworkRetriever(
		c,
		config.ConfigLoadRoots{
			CWDRoot: os.DirFS(cwd),
			Home:    os.DirFS(homeDir),
		},
	)
	if err != nil {
		return err
	}

	for _, getter := range frameworkGetters {
		if err := frameworkRetriever.Get(getter); err != nil {
			return err
		}
	}

	return nil
}
