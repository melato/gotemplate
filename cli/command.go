/*
Package cli implement a Command Line Interface, using melato.org/command.
*/
package cli

import (
	_ "embed"

	"melato.org/command"
	"melato.org/command/usage"
	"melato.org/gotemplate"
)

//go:embed usage.yaml
var usageData []byte

func Command(t *gotemplate.TemplateOp) *command.SimpleCommand {
	var cmd command.SimpleCommand
	cmd.Command("exec").Flags(t).RunFunc(t.Run)
	help := cmd.Command("help")
	help.Command("templates").RunFunc(t.ListTemplates)
	help.Command("funcs").RunFunc(t.ListFuncs)
	help.Command("globals").RunFunc(t.ListGlobals)
	help.Command("formats").RunFunc(t.ListPropertyFormats)
	help.Command("func").RunFunc(t.FuncUsage)
	usage.Apply(&cmd, usageData)
	return &cmd
}
