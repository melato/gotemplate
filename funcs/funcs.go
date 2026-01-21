package funcs

import (
	_ "embed"
)

//go:embed funcs.txt
var funcUsage []byte

type Funcs interface {
	SetFunc(string, any)
	AddUsageTxt([]byte)
}

func AddFuncs(funcs Funcs) {
	funcs.SetFunc("file", ReadFile)
	funcs.AddUsageTxt(funcUsage)
}
