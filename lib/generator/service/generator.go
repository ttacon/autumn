package service

import (
	"bytes"
	"errors"
	"html/template"
	"os"

	"github.com/ttacon/autumn/lib/config"
	"github.com/ttacon/autumn/lib/engine"
)

// ServiceGenerator generates the service file content for a given model.
type ServiceGenerator interface {
	GenerateContent(m engine.ModelTarget) ([]byte, error)
	GenerateAndStoreContent(m engine.ModelTarget, path string) error
}

type serviceGenerator struct {
	framework           config.Framework
	templatesToGenerate []string
}

func NewServiceGenerator(
	frameworkName string,
	frameworkSource config.FrameworkSource,
	templatesToGenerate []string,
) (ServiceGenerator, error) {
	framework, exists := frameworkSource.GetFramework(frameworkName)
	if !exists {
		return nil, ErrNoSuchFramework
	}

	// If not specific templates are provided, default to the core templates.
	if templatesToGenerate == nil {
		templatesToGenerate = DefaultTemplates
	}

	return &serviceGenerator{
		framework:           framework,
		templatesToGenerate: templatesToGenerate,
	}, nil
}

var (
	ErrNoSuchFramework = errors.New("no such framework exists")
	ErrNoSuchTemplate  = errors.New("no such template")
)

var DefaultTemplates = []string{
	"CreateTemplate",
	"RetrieveTemplate",
	"UpdateTemplate",
	"DeleteTemplate",
	"ListTemplate",
	// NOTE(ttacon): should we support a SearchTemplate?
}

func (sg *serviceGenerator) GenerateContent(m engine.ModelTarget) ([]byte, error) {

	var buf = bytes.NewBuffer(nil)

	for _, templName := range sg.templatesToGenerate {
		templRaw, ok := sg.framework.GetTemplate(templName)
		if !ok {
			return nil, ErrNoSuchTemplate
		}

		templ, err := template.
			New("service generation template: " + templName).
			Parse(string(templRaw))
		if err != nil {
			return nil, err
		} else if err := templ.Execute(buf, m.ToTemplateVariables()); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (sg *serviceGenerator) hydrateTemplateVariables(
	tmplVars map[string]interface{},
) map[string]interface{} {

	// NOTE(ttacon): we'll want to source this from the generator config in the future
	// and default back to "services".
	tmplVars["PackageName"] = "services"

	return tmplVars
}

func (sg *serviceGenerator) GenerateAndStoreContent(
	m engine.ModelTarget,
	path string,
) error {
	data, err := sg.GenerateContent(m)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
