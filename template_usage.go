package gotemplate

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func (t *TemplateOp) GetUsage() (map[string]string, error) {
	if !t.parsedUsage {
		t.funcUsage = make(map[string]string)
		t.parsedUsage = true
		for _, data := range t.funcUsageTxt {
			ParseUsage(data, "", t.funcUsage)
		}
	}
	return t.funcUsage, nil
}

func (t *TemplateOp) ListFuncs() error {
	usage, err := t.GetUsage()
	if err != nil {
		return err
	}
	names := make([]string, 0, len(t.Funcs))
	for name, _ := range t.Funcs {
		names = append(names, name)
	}
	sort.Strings(names)
	maxlen := maxRunes(names)
	for _, name := range names {
		summary := ""
		u, found := usage[name]
		if found {
			summary = firstLine(u)
		}
		fmt.Printf("%-*s %s\n", maxlen, name, summary)
	}
	return nil
}

func (t *TemplateOp) funcSignature(name string, fType reflect.Type, isMethod bool) {
	n := fType.NumIn()
	offset := 0
	if isMethod {
		if n == 0 {
			return
		}
		offset = 1
		n -= 1
	}
	args := make([]string, n)
	for i := 0; i < n; i++ {
		pType := fType.In(offset + i)
		args[i] = fmt.Sprintf("%v", pType)
	}
	fmt.Printf("%s(%s)\n", name, strings.Join(args, ", "))
}

func (t *TemplateOp) fUsage(name string, fType reflect.Type) error {
	usage, err := t.GetUsage()
	if err != nil {
		return err
	}
	desc, found := usage[name]
	t.funcSignature(name, fType, false)
	if found {
		for line := range iterLines(strings.NewReader(desc)) {
			fmt.Printf("   %s\n", line)
		}
	} else if fType.NumIn() == 0 && fType.NumOut() > 0 {
		outType := fType.Out(0)
		n := outType.NumMethod()
		if n > 0 {
			fmt.Printf("methods:\n")
			for i := 0; i < n; i++ {
				method := outType.Method(i)
				t.funcSignature(method.Name, method.Type, true)
			}
		}
	}
	return nil
}

func (t *TemplateOp) FuncUsage(name string) error {
	f, found := t.Funcs[name]
	if found {
		fType := reflect.TypeOf(f)
		return t.fUsage(name, fType)
	}
	return fmt.Errorf("no such func: %s", name)
}

func (t *TemplateOp) ListGlobals() error {
	globals := parseGlobal()
	names := make([]string, 0, len(globals))
	for name, _ := range globals {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\n", name)
		for line := range iterLines(strings.NewReader(globals[name])) {
			fmt.Printf("     %s\n", line)
		}
	}
	return nil
}

func (t *TemplateOp) ListPropertyFormats() {
	names := make([]string, 0, len(propertyParsers))
	for name, _ := range propertyParsers {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s\n", name)
	}
}
