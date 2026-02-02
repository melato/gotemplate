package gotemplate

import (
	"fmt"
	"reflect"
	"text/template"
)

var TraceFuncs bool

// Funcs is used to create a template.FuncMap
// If a value is a function, it is copied as is to the FuncMap.
// Otherwise, it is converted to a func that returns it.
type Funcs map[string]any

/*
Generate a template.FuncMap
*/
func (t Funcs) CreateFuncMap() template.FuncMap {
	fm := make(template.FuncMap)
	for name, v := range t {
		if TraceFuncs {
			fmt.Printf("func %s: %T\n", name, v)
		}
		if reflect.TypeOf(v).Kind() == reflect.Func {
			fm[name] = v
			continue
		}
		// create a func that returns v
		fm[name] = func() any { return v }
	}
	return fm
}
