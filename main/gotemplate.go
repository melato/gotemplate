package main

import (
	_ "embed"
	"fmt"

	"melato.org/command"
	"melato.org/gotemplate"
	"melato.org/gotemplate/funcs"
)

var version string

//go:embed funcs.yaml
var funcUsage []byte

func DefineFuncs(op *gotemplate.TemplateOp) {
	op.SetFunc("file", funcs.ReadFile)
	op.AddUsageYaml(funcUsage)
}

func main() {
	var op gotemplate.TemplateOp
	DefineFuncs(&op)
	cmd := op.Command()
	cmd.Command("version").NoConfig().RunFunc(func() {
		fmt.Printf("%s\n", version)
	}).Short("print version")
	command.Main(cmd)
}
