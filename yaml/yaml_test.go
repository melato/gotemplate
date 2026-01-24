package yaml

import (
	_ "embed"
	"testing"
)

//go:embed yaml_test.yaml
var yamlData []byte

func TestParseYaml(t *testing.T) {
	m := make(map[string]any)
	m["c"] = "C"
	err := ParseYaml(yamlData, m)
	if err != nil {
		t.Fatalf("%v", err)
	}
	verify := func(key, value string) {
		x := m[key]
		if value != x {
			t.Fatalf("%s: %v != %s", key, x, value)
		}
	}
	verify("a", "A")
	verify("b", "B")
	verify("c", "C")
}
