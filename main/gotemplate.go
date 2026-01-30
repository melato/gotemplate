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
	var config gotemplate.Config
	funcs.AddFuncs(&config)
	cmd := cli.Command(&config)
	cmd.Command("version").NoConfig().RunFunc(func() {
		fmt.Printf("%s\n", version)
	}).Short("print version")
	command.Main(cmd)
}
