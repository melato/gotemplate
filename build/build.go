package build

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
	"melato.org/gotemplate"
)

type BuildOp struct {
	Funcs gotemplate.Funcs `name:"-"`

	ConfigFile string `name:"c" usage:"build config file"`
	InputDir   string `name:"i" usage:"input dir, overrides config input_dir"`
	OutputDir  string `name:"o" usage:"output dir, overrides config output_dir"`
	Verbose    bool   `name:"v" usage:"verbose, print template names"`
}

func (t *BuildOp) Build(args ...string) error {
	if t.ConfigFile == "" {
		return fmt.Errorf("missing config file")
	}
	configDir := filepath.Dir(t.ConfigFile)
	funcs, err := t.Funcs.CreateFuncMap(configDir)
	if err != nil {
		return err
	}
	var config Config
	data, err := os.ReadFile(t.ConfigFile)
	if err == nil {
		err = yaml.Unmarshal(data, &config)
	}
	if err != nil {
		return err
	}
	commonTpl := template.New("")
	commonTpl.Funcs(funcs)
	for _, tc := range config.Templates {
		dir := ResolvePath(configDir, tc.Dir)
		if t.Verbose {
			fmt.Printf("using templates from %s.  patterns: %v\n", dir, tc.Patterns)
		}
		f := os.DirFS(dir)
		_, err = commonTpl.ParseFS(f, tc.Patterns...)
		if err != nil {
			return err
		}
	}
	model := make(map[any]any)
	for name, value := range config.Properties {
		model[name] = value
	}
	inputDir := t.InputDir
	if inputDir == "" {
		inputDir = ResolvePath(configDir, config.InputDir)
	}
	outputDir := t.OutputDir
	if outputDir == "" {
		outputDir = ResolvePath(configDir, config.OutputDir)
	}
	err = os.MkdirAll(outputDir, os.FileMode(0755))
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(inputDir)
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
		inputFile := filepath.Join(inputDir, name)
		outputFile := filepath.Join(outputDir, strings.TrimSuffix(name, ext)+config.OutputExtension)

		tpl, err := commonTpl.Clone()
		if err != nil {
			return err
		}
		if t.Verbose {
			fmt.Printf("parsing %s\n", inputFile)
		}
		_, err = tpl.ParseFiles(inputFile)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if t.Verbose {
			fmt.Printf("execute %s\n", name)
		}
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
