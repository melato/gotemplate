package gotemplate

import (
	"io/fs"
)

type TemplateSet struct {
	FS       fs.FS
	Patterns []string
}
