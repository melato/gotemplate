package main

import (
	"fmt"

	"melato.org/command"
	"melato.org/gotemplate"
	"melato.org/gotemplate/cli"
	"melato.org/gotemplate/funcs"
)

var version string

func main() {
	var op gotemplate.TemplateOp
	funcs.AddFuncs(&op)
	cmd := cli.Command(&op)
	cmd.Command("version").NoConfig().RunFunc(func() {
		fmt.Printf("%s\n", version)
	}).Short("print version")
	command.Main(cmd)
}
