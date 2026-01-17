package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/gotemplate"
)

var version string

//go:embed usage.yaml
var usageData []byte

type App struct {
	gotemplate.TemplateOp
	FileFunctions bool `name:"files" usage:"define functions that read files"`
	Version       bool `name:"version" usage:"print version and exit"`
}

func (t *App) DefineFuncs() {
	if !t.FileFunctions {
		return
	}
	var f gotemplate.FileFunctions
	t.SetFunc("file", f.File)
	t.SetFunc("json", f.Json)
	t.SetFunc("yaml", f.Yaml)
}

func (t *App) Configured() error {
	err := t.TemplateOp.Configured()
	if err != nil {
		return err
	}
	t.DefineFuncs()
	return nil
}

func (t *App) Run(args ...string) error {
	if t.Version {
		fmt.Printf("%s\n", version)
		return nil
	}
	return t.TemplateOp.Run(args)
}

func main() {
	var cmd command.SimpleCommand
	var app App
	cmd.Flags(&app).RunFunc(app.Run)
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
