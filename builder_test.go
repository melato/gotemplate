package gotemplate

import (
	"testing"
)

func TestBuilder(t *testing.T) {
	var builder builder
	builder.values = make(map[string]any)
	builder.Set("a", "1")
	builder.Set("b", "2")
	verify := func(key, value string) {
		v := builder.values[key]
		if v != value {
			t.Fatalf("%s=%v != %s", key, v, value)
		}
	}
	verify("a", "1")
	verify("b", "2")
	builder.Set("", map[string]any{"c": "3"})
	verify("c", "3")
	verify("a", "1")
}
