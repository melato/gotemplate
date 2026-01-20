package gotemplate

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

func (t *TemplateOp) GetUsage() (map[string]FuncUsage, error) {
	if !t.parsedUsage {
		t.parsedUsage = true
		for _, data := range t.funcUsageYaml {
			if t.funcUsage == nil {
				t.funcUsage = make(map[string]FuncUsage)
			}
			err := yaml.Unmarshal(data, &t.funcUsage)
			if err != nil {
				return nil, err
			}
		}
	}
	return t.funcUsage, nil
}

func (t *TemplateOp) ListFuncs() error {
	usage, err := t.GetUsage()
	if err != nil {
		return err
	}
	// compute the number of runes in a string
	runeCount := func(s string) int {
		var i int
		for i, _ = range s {
		}
		return i + 1
	}
	var maxlen int
	names := make([]string, 0, len(t.Funcs))
	for name, _ := range t.Funcs {
		names = append(names, name)
		w := runeCount(name)
		if w > maxlen {
			maxlen = w
		}
	}
	sort.Strings(names)
	for _, name := range names {
		summary := ""
		u, found := usage[name]
		if found {
			summary = firstLine(u.Description)
		}
		fmt.Printf("%*s: %s\n", maxlen, name, summary)
	}
	return nil
}

func (t *TemplateOp) fUsage(name string, fType reflect.Type) error {
	usage, err := t.GetUsage()
	if err != nil {
		return err
	}
	u, found := usage[name]
	n := fType.NumIn()
	params := make([]string, n)
	for i := 0; i < n; i++ {
		pType := fType.In(i)
		if found && n == len(u.Params) {
			params[i] = fmt.Sprintf("%s %v", u.Params[i].Name, pType)
		} else {
			params[i] = fmt.Sprintf("%v", pType)
		}
	}
	fmt.Printf("%s(%s)\n", name, strings.Join(params, ", "))
	if found {
		fmt.Printf("%s\n", strings.TrimSpace(u.Description))
		if len(u.Params) > 0 {
			fmt.Printf("\nParameters:\n")
		}
		for _, param := range u.Params {
			fmt.Printf("%s: %s\n", param.Name, param.Description)
		}
	} else if n == 0 && fType.NumOut() > 0 {
		outType := fType.Out(0)
		n := outType.NumMethod()
		if n > 0 {
			fmt.Printf("methods:\n")
			for i := 0; i < n; i++ {
				fmt.Printf(" %s\n", outType.Method(i).Name)
			}
		}
	}
	return nil
}

func (t *TemplateOp) FuncUsage(name string) error {
	f, found := t.Funcs[name]
	if !found {
		return fmt.Errorf("no such func: %s", name)
	}
	fType := reflect.TypeOf(f)
	return t.fUsage(name, fType)
}
