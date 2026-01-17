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
	TemplateFile string `name:"t" usage:"template file"`
	OutputFile   string `name:"o" usage:"output file"`
	FileMode     string `name:"mode" usage:"output file mode"`
	Delims       string `name:"delims" usage:"template left,right delims, separated by ','"`
	leftDelim    string
	rightDelim   string

	Options
	Funcs template.FuncMap `name:"-"`
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

func (t *TemplateOp) buildTemplate() (*template.Template, error) {
	var data []byte
	var err error
	if t.TemplateFile != "" {
		data, err = t.readFile(t.TemplateFile)
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

	tpl := template.New("x")
	if t.leftDelim != "" || t.rightDelim != "" {
		tpl.Delims(t.leftDelim, t.rightDelim)
	}
	tpl.Funcs(t.Funcs)
	_, err = tpl.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t *TemplateOp) Run() error {
	tpl, err := t.buildTemplate()
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
