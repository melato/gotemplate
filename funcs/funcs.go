package funcs

import (
	_ "embed"
)

//go:embed funcs.yaml
var funcUsage []byte

type Funcs interface {
	SetFunc(string, any)
	AddUsageYaml([]byte)
}

func AddFuncs(funcs Funcs) {
	funcs.SetFunc("file", ReadFile)
	funcs.AddUsageYaml(funcUsage)
}
