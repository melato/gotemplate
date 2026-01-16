package gotemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
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
	t.Funcs = make(template.FuncMap)
	t.Funcs["file"] = func(file string) (string, error) {
		data, err := t.readFile(file)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	t.Funcs["json"] = func(file string) (any, error) {
		data, err := t.readFile(file)
		if err != nil {
			return nil, err
		}
		var v any
		err = json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	t.Funcs["yaml"] = func(file string) (map[string]any, error) {
		data, err := t.readFile(file)
		if err != nil {
			return nil, err
		}
		var v map[string]any
		err = yaml.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil
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

func (t *TemplateOp) buildTemplate(data []byte) (*template.Template, error) {
	tpl := template.New("x")
	if t.leftDelim != "" || t.rightDelim != "" {
		tpl.Delims(t.leftDelim, t.rightDelim)
	}
	tpl.Funcs(t.Funcs)
	_, err := tpl.Parse(string(data))
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (t *TemplateOp) Run() error {
	values, err := t.Values()
	if err != nil {
		return err
	}
	var data []byte
	if t.TemplateFile != "" {
		data, err = t.readFile(t.TemplateFile)
	} else {
		var buf bytes.Buffer
		_, err = io.Copy(&buf, os.Stdin)
		if err == nil {
			data = buf.Bytes()
		}
	}
	tpl, err := t.buildTemplate(data)
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
