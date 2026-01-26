package gotemplate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

type TemplateOp struct {
	Funcs        template.FuncMap
	TemplateName string `name:"t" usage:"template name or file"`
	OutputFile   string `name:"o" usage:"output file"`
	FileMode     string `name:"mode" usage:"output file mode"`
	Delims       string `name:"delims" usage:"template left,right delims, separated by ','"`
	leftDelim    string
	rightDelim   string

	Options
}

func (t *TemplateOp) Init() error {
	t.FileMode = "0664"
	t.Delims = "{{,}}"
	return nil
}

func (t *TemplateOp) SetFunc(name string, f any) {
	if t.Funcs == nil {
		t.Funcs = make(template.FuncMap)
	}
	t.Funcs[name] = f
}

func (t *TemplateOp) Configured() error {
	if t.Delims != "" {
		left, right, ok := strings.Cut(t.Delims, ",")
		if !ok {
			return fmt.Errorf("invalid delims: %s", t.Delims)
		}
		t.leftDelim = left
		t.rightDelim = right
	}
	return nil
}

func (t *TemplateOp) newTemplate() *template.Template {
	tpl := template.New("")
	if t.leftDelim != "" || t.rightDelim != "" {
		tpl.Delims(t.leftDelim, t.rightDelim)
	}
	tpl.Funcs(t.Funcs)
	return tpl
}

func (t *TemplateOp) buildSingleTemplate() (*template.Template, error) {
	var data []byte
	var err error
	if t.TemplateName != "" {
		data, err = t.readFile(t.TemplateName)
	} else {
		var buf bytes.Buffer
		_, err = io.Copy(&buf, os.Stdin)
		if err == nil {
			data = buf.Bytes()
		}
	}
	if err != nil {
		return nil, err
	}

	tpl := t.newTemplate()
	_, err = tpl.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t *TemplateOp) buildTemplate(args []string) (*template.Template, error) {
	if len(args) == 0 {
		return t.buildSingleTemplate()
	}
	tpl := t.newTemplate()
	tpl, err := tpl.ParseFiles(args...)
	if err != nil {
		return nil, err
	}
	if t.TemplateName != "" {
		t0 := tpl.Lookup(t.TemplateName)
		if t0 == nil {
			return nil, fmt.Errorf("no such template: %s", t.TemplateName)
		}
		return t0, nil
	} else {
		for _, t0 := range tpl.Templates() {
			return t0, nil
		}
		return nil, fmt.Errorf("no templates")
	}
	return tpl, nil
}

func listTemplates(tpl *template.Template) {
	fmt.Printf("Templates:\n")
	for _, tpl := range tpl.Templates() {
		fmt.Printf("  %s\n", tpl.Name())
	}
}

func (t *TemplateOp) ListTemplates(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing template files")
	}
	tpl := t.newTemplate()
	tpl, err := tpl.ParseFiles(args...)
	if err != nil {
		return err
	}
	listTemplates(tpl)
	return nil
}

func (t *TemplateOp) Run(args ...string) error {
	tpl, err := t.buildTemplate(args)
	if err != nil {
		return err
	}
	values, err := t.Values()
	if err != nil {
		return err
	}
	w := os.Stdout
	var tmpName string
	if t.OutputFile != "" {
		mode, err := strconv.ParseInt(t.FileMode, 8, 32)
		if err != nil {
			return fmt.Errorf("invalid mode: %s", t.FileMode)
		}
		dir := filepath.Dir(t.OutputFile)
		w, err = os.CreateTemp(dir, "tpl")
		if err != nil {
			return nil
		}
		tmpName = w.Name()
		defer os.Remove(tmpName)
		err = os.Chmod(tmpName, os.FileMode(mode))
		if err != nil {
			return nil
		}
		w, err = os.OpenFile(tmpName, os.O_RDWR, os.FileMode(mode))
		if err != nil {
			return err
		}
		defer w.Close()
	}
	err = tpl.Execute(w, values)
	if t.OutputFile != "" {
		w.Close()
		if err != nil {
			return err
		}
		err = os.Rename(tmpName, t.OutputFile)
	}
	return err
}
