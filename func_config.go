package gotemplate

import (
	"text/template"
)

type FuncConfig struct {
	Funcs        template.FuncMap
	funcUsage    map[string]string
	funcUsageTxt [][]byte
	parsedUsage  bool
}

func (t *FuncConfig) SetFunc(name string, f any) {
	if t.Funcs == nil {
		t.Funcs = make(template.FuncMap)
	}
	t.Funcs[name] = f
}

/*
Add usage for functions
*/
func (t *FuncConfig) AddUsage(funcUsage map[string]string) {
	if t.funcUsage == nil {
		t.funcUsage = make(map[string]string)
	}
	for name, u := range funcUsage {
		t.funcUsage[name] = u
	}
}

/*
Add usage for functions, in text format
The usage is parsed when needed
*/
func (t *FuncConfig) AddUsageTxt(usageTxt []byte) {
	t.funcUsageTxt = append(t.funcUsageTxt, usageTxt)
}
