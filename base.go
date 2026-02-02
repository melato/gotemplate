package gotemplate

import (
	"io/fs"
	"text/template"
)

type BaseConfig struct {
	Funcs     Funcs
	Templates []TemplateSet
}

type TemplateSet struct {
	FS       fs.FS
	Patterns []string
}

func (t *BaseConfig) Apply(tpl *template.Template) error {
	funcs := t.Funcs.CreateFuncMap()
	tpl.Funcs(funcs)
	for _, tc := range t.Templates {
		_, err := tpl.ParseFS(tc.FS, tc.Patterns...)
		if err != nil {
			return err
		}
	}
	return nil
}
