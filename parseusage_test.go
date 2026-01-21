package gotemplate

import (
	_ "embed"
	"testing"
)

//go:embed test/doc.txt
var testUsage []byte

func TestParseUsage(t *testing.T) {
	usage := make(map[string]string)
	ParseUsage(testUsage, "", usage)
	verify := func(name, desc string) {
		s := usage[name]
		if s != desc {
			t.Fatalf("%s: %q != %q\n", name, s, desc)
		}
	}
	verify("a", "a1")
	verify("b", "b1\n b2")
}

func TestGlobal(t *testing.T) {
	usage := make(map[string]string)
	ParseUsage(globalUsage, "", usage)
	for _, name := range []string{"and"} {
		_, found := usage[name]
		if !found {
			t.Fatalf("global not found: %s", name)
		}
	}
}
