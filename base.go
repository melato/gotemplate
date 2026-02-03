package gotemplate

import (
	"io/fs"
	"text/template"
)

type BaseConfig struct {
	Funcs      Funcs
	Templates  []TemplateSet
	Properties map[string]any
}

type TemplateSet struct {
	FS       fs.FS
	Patterns []string
}

func (t *BaseConfig) SetFunc(name string, f any) {
	if t.Funcs == nil {
		t.Funcs = make(map[string]any)
	}
	t.Funcs[name] = f
}

func (t *BaseConfig) SetProperty(name string, value any) {
	if t.Properties == nil {
		t.Properties = make(map[string]any)
	}
	t.Properties[name] = value
}

func (t *BaseConfig) Apply(tpl *template.Template) error {
	funcs := t.Funcs.CreateFuncMap()
	tpl.Funcs(funcs)
	tpl.Option("missingkey=error")
	for _, tc := range t.Templates {
		_, err := tpl.ParseFS(tc.FS, tc.Patterns...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *BaseConfig) CreateModel() map[any]any {
	model := make(map[any]any)
	for name, value := range t.Properties {
		model[name] = value
	}
	return model
}
