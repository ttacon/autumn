package engine

import (
	"testing"
	"testing/fstest"
	"time"
)

// Test files
var (
	modelGoFile = `
package models // github.com/ttacon/example-foo/models

// ResourceModel is a resource model that we want to generate a service and
// controller for.
//
// @Autumn:Model
type ResourceModel struct {
    ID string
    Name string
    Email string
}`
	textFile = `this is a text file`
)

func TestNewEngine(t *testing.T) {
	var rootFS = fstest.MapFS{
		"root/model.go": &fstest.MapFile{
			Data:    []byte(modelGoFile),
			Mode:    0644,
			ModTime: time.Now(),
			Sys:     nil,
		},
		"root/README": &fstest.MapFile{
			Data:    []byte(textFile),
			Mode:    0644,
			ModTime: time.Now(),
			Sys:     nil,
		},
	}

	eng, err := NewEngine(rootFS)
	if err != nil {
		t.Error("unexpected err: ", err)
		t.Fail()
	}

	modelTargets, err := eng.IdentifyModelTargets()
	if err != nil {
		t.Error("unexpected err: ", err)
	} else if len(modelTargets) != 1 {
		t.Error("expected 1 Go file to be found, found: ", len(modelTargets))
	}

	model := modelTargets[0]
	if modl, ok := model.(*modelTarget); !ok {
		t.Error("expect modl to be a modelTarget, somehow was not")
	} else if modl.pkgName != "models" {
		t.Error("expected pkgName to be 'models', found: " + modl.pkgName)
	}

	if modelName, err := model.Name(); err != nil {
		t.Error("failed to get model name: ", err)
	} else if modelName != "ResourceModel" {
		t.Error("got unexpected model name: ", modelName)
	}
}
