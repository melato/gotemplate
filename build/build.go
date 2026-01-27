package build

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type BuildOp struct {
	Funcs template.FuncMap `name:"-"`

	ConfigFile string `name:"c" usage:"build config file"`
}

func (t *BuildOp) Build(args ...string) error {
	var config Config
	data, err := os.ReadFile(t.ConfigFile)
	if err == nil {
		err = yaml.Unmarshal(data, &config)
	}
	if err != nil {
		return err
	}
	commonTpl := template.New("")
	commonTpl.Funcs(t.Funcs)
	f := os.DirFS(config.Template.Dir)
	_, err = commonTpl.ParseFS(f, config.Template.Patterns...)
	if err != nil {
		return err
	}
	model := make(map[any]any)
	for name, value := range config.Properties {
		model[name] = value
	}
	entries, err := os.ReadDir(config.InputDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != config.InputExtension {
			continue
		}
		inputFile := filepath.Join(config.InputDir, name)
		outputFile := filepath.Join(config.OutputDir, strings.TrimSuffix(name, ext)+config.OutputExtension)

		tpl, err := commonTpl.Clone()
		if err != nil {
			return err
		}
		_, err = tpl.ParseFiles(inputFile)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = tpl.ExecuteTemplate(&buf, name, model)
		if err != nil {
			return fmt.Errorf("ExecuteTemplate: %w", err)
		}
		err = os.WriteFile(outputFile, buf.Bytes(), os.FileMode(0664))
		if err != nil {
			return err
		}
	}
	return nil
}
