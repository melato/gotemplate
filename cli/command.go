/*
Package cli implement a Command Line Interface, using melato.org/command.
*/
package cli

import (
	_ "embed"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/gotemplate"
	"melato.org/gotemplate/build"
)

//go:embed usage.yaml
var usageData []byte

func Command(fc *gotemplate.Config) *command.SimpleCommand {
	var cmd command.SimpleCommand
	var op gotemplate.TemplateOp
	op.Base = fc.BaseConfig
	cmd.Command("exec").Flags(&op).RunFunc(op.Run)

	var buildOp build.BuildOp
	buildOp.Base = fc.BaseConfig
	cmd.Command("build").Flags(&buildOp).RunFunc(buildOp.Build)

	help := cmd.Command("help")
	help.Command("templates").RunFunc(op.ListTemplates)
	help.Command("funcs").RunFunc(fc.ListFuncs)
	help.Command("globals").RunFunc(gotemplate.ListGlobals)
	help.Command("formats").RunFunc(gotemplate.ListPropertyFormats)
	help.Command("func").RunFunc(fc.FuncUsage)
	usage.Apply(&cmd, usageData)
	return &cmd
}
