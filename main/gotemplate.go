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
	Version bool `name:"version" usage:"print version and exit"`
}

func (t *App) Run() error {
	if t.Version {
		fmt.Printf("%s\n", version)
		return nil
	}
	return t.TemplateOp.Run()
}

func main() {
	var cmd command.SimpleCommand
	var app App
	cmd.Flags(&app).RunFunc(app.Run)
	usage.Apply(&cmd, usageData)
	command.Main(&cmd)
}
