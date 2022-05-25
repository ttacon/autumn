package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ttacon/autumn/lib/config"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/mod/modfile"
)

func initCommand(c *cli.Context) error {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("failed to identify current working directory: ", err)
		return err
	}
	root := os.DirFS(cwd)

	if err := ensureAutumnDirExists(); err != nil {
		fmt.Println("failed to ensure that the .autumn directory exists")
		return err
	}

	if configExists, err := checkForConfigFile(root); err != nil {
		fmt.Println("failed to check for existing config file: ", err)
		return err
	} else if configExists {
		if !c.Bool("force") {
			fmt.Println("found existing configuration file, exiting")
			return nil
		}
		fmt.Println("existing config file found, but force was used, continuing")
	}

	var conf config.Config

	packageName, err := pkgNameFromModFile(root)
	if err != nil {
		fmt.Println("failed to idenfity the package name from the mod file, err: ", err)
		return err
	} else if len(packageName) == 0 {
		fmt.Println("no package name found from mod file")
		return errors.New("mod file expected, none found")
	}
	conf.Name = packageName

	configFilePath := filepath.Join(
		autumnDir,
		configFileName,
	)

	var buf = bytes.NewBuffer(nil)
	if err := toml.NewEncoder(buf).Encode(conf); err != nil {
		fmt.Println("failed to prepare config for storage: ", err)
		return err
	} else if err := ioutil.WriteFile(configFilePath, buf.Bytes(), 0644); err != nil {
		fmt.Println("failed to write file to storage: ", err)
		return err
	}

	return retrieveSourcesForEngine(conf)
}

var (
	configFileName = "config"
)

func checkForConfigFile(root fs.FS) (bool, error) {
	if _, err := root.Open(
		filepath.Join(
			autumnDir,
			configFileName,
		),
	); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func pkgNameFromModFile(root fs.FS) (string, error) {
	file, err := root.Open("go.mod")
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	modFile, err := modfile.ParseLax("go.mod", data, nil)
	if err != nil {
		return "", err
	}

	return modFile.Module.Mod.Path, nil
}

var (
	autumnDir = ".autumn"
)

func ensureAutumnDirExists() error {
	return os.MkdirAll(
		filepath.Join(
			autumnDir,
			"frameworks",
		),
		os.ModeDir|0755,
	)
}
