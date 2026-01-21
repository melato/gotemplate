package gotemplate

import (
	_ "embed"
)

//go:embed global/template.txt
var globalUsage []byte

func parseGlobal() map[string]string {
	globals := make(map[string]string)
	ParseUsage(globalUsage, "Predefined", globals)
	return globals
}
