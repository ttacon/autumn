package config

type FrameworkSource interface {
	AddFramework(name string, f Framework) FrameworkSource
	GetFramework(name string) (Framework, bool)
}

type frameworkSource map[string]Framework

func (fs frameworkSource) GetFramework(name string) (Framework, bool) {
	framework, ok := fs[name]
	return framework, ok
}

type Framework interface {
	AddTemplate(name string, data []byte) Framework
	GetTemplate(name string) ([]byte, bool)
}

type framework map[string][]byte

func (f framework) GetTemplate(name string) ([]byte, bool) {
	data, ok := f[name]
	return data, ok
}

func NewFrameworkSource() FrameworkSource {
	return make(frameworkSource)
}

func (fs frameworkSource) AddFramework(name string, f Framework) FrameworkSource {
	fs[name] = f
	return fs
}

func (f framework) AddTemplate(name string, data []byte) Framework {
	f[name] = data
	return f
}

func NewFramework() Framework {
	return make(framework)
}

func FrameworkSourceFromMap(data map[string]map[string][]byte) FrameworkSource {
	fs := NewFrameworkSource()
	for name, templates := range data {
		f := NewFramework()
		for templateName, data := range templates {
			f.AddTemplate(templateName, data)
		}
		fs.AddFramework(name, f)
	}
	return fs
}
