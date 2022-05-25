package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ttacon/autumn/lib/config"
	"github.com/ttacon/autumn/lib/engine"
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
	if _, err := root.Open(outputFileName); err == nil && !c.Bool("force") {
		return errors.New("output file already exists")
	}

	// NOTE(ttacon): punting on targeting mode for now since we support
	// other methods for specifying generation targts.

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

	// NOTE(ttacon): Instead of retrieving assets, should we report missing
	// dependencies and support retrieveing them now via a flag?
	if err := retrieveSourcesForEngine(conf); err != nil {
		return err
	}

	// Load in the engine.
	eng, err := engine.NewEngine(root)
	if err != nil {
		return err
	}

	// Identify our targets to create.
	targets, err := eng.IdentifyModelTargets()
	if err != nil {
		return err
	}

	// Generate plan (Model -> <Controller, Router, Service> -> Framework -> Templates).
	var planData = make(map[string]PlanData)
	for _, target := range targets {
		name, err := target.Name()
		if err != nil {
			return err
		}
		planData[name] = PlanData{
			Config: conf,
			Model: ModelTargetPlan{
				Name:        name,
				PackageName: target.PkgName(),
				Raw:         target,
			},
		}
	}

	rawPlanData, err := json.Marshal(planData)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(outputFileName, rawPlanData, 0644)
}

type PlanData struct {
	Model  ModelTargetPlan
	Config config.Config
}

type ModelTargetPlan struct {
	Name        string
	PackageName string
	Raw         interface{} // This should be versioned
}
