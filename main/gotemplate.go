package main

import (
	"fmt"

	"melato.org/command"
	"melato.org/gotemplate"
	"melato.org/gotemplate/cli"
	"melato.org/gotemplate/funcs"
	"melato.org/gotemplate/yaml"
)

var version string

func main() {
	gotemplate.SetParser("yaml", yaml.ParseYaml)
	var config gotemplate.Config
	funcs.AddFuncs(&config)
	cmd := cli.Command(&config)
	cmd.Command("version").NoConfig().RunFunc(func() {
		fmt.Printf("%s\n", version)
	}).Short("print version")
	command.Main(cmd)
}
