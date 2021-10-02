package service

import (
	"testing"
	"testing/fstest"
	"time"

	"github.com/ttacon/autumn/lib/config"
	"github.com/ttacon/autumn/lib/engine"
)

// Test files
var (
	modelGoFile = `
package models // github.com/ttacon/example-foo/models

// ResourceModel is a resource model that we want to generate a service and
// controller for.
//
// @Autumn:Model
type ResourceFoo struct {
    ID string
    Name string
    Email string
}`
	textFile = `this is a text file`
)

func TestNewServiceGenerator(t *testing.T) {
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

	eng, err := engine.NewEngine(rootFS)
	if err != nil {
		t.Error("unexpected err: ", err)
		t.Fail()
	}

	modelTargets, err := eng.IdentifyModelTargets()
	if err != nil {
		t.Error("unexpected err: ", err)
		t.Fail()
	}

	if len(modelTargets) > 1 {
		t.Error("expected only one target")
		t.Fail()
	}

	model := modelTargets[0]

	fs := config.FrameworkSourceFromMap(map[string]map[string][]byte{
		"go.mongodb.org/mongo-driver/mongo": map[string][]byte{
			"CreateTemplate": []byte(`func Create{{.Name}}(){  }`),
		},
	})

	gener8r, err := NewServiceGenerator(
		"go.mongodb.org/mongo-driver/mongo",
		fs,
		[]string{"CreateTemplate"},
	)
	if err != nil {
		t.Error("unexpected err: ", err)
		t.Fail()
	}

	expectedFile := `func CreateResourceFoo(){  }`

	data, err := gener8r.GenerateContent(model)
	if err != nil {
		t.Error("unexpected err: ", err)
	} else if string(data) != expectedFile {
		t.Error("received unexpected file content: ", string(data))
	}

}
