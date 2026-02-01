package gotemplate

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"text/template"
)

// Funcs is used to create a template.FuncMap
// If a value is a function, it is copied as is to the FuncMap.
// If a value implements FuncFC, its Func() method is called
// along with an appropriate fs.FS,
// and the return value (which should be a function),
// is placed in the FuncMap.
// If the value is not any of the above types,
// it is converted to a func that returns it.
type Funcs map[string]any

// FuncFC allows using template functions that have access to an fs.FS filesystem,
// that they can use to read files from.
type FuncFC interface {
	Func(moduleFS fs.FS) (any, error)
}

/*
Generate a template.FuncMap, using os.DirFS(configDir),
to generate functions from FuncFC.

If a configuration file is used to configure the template,
configDir should be the directory of the configuration file.
Otherwise, it should be the current directory.
*/
func (t Funcs) CreateFuncMap(configDir string) (template.FuncMap, error) {
	configFS := os.DirFS(configDir)
	fm := make(template.FuncMap)
	for name, v := range t {
		if reflect.TypeOf(v).Kind() == reflect.Func {
			fm[name] = v
			continue
		}
		fc, ok := v.(FuncFC)
		if ok {
			f, err := fc.Func(configFS)
			if err != nil {
				return nil, fmt.Errorf("func %s.Func(FS) error: %w", name, err)
			}
			fm[name] = f
			continue
		}
		// create a func that returns v
		fm[name] = func() any { return v }
	}
	return fm, nil
}
