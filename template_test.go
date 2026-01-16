package gotemplate

import (
	"bytes"
	"embed"
	"testing"
)

//go:embed test/*
var testFS embed.FS

func TestOptions(t *testing.T) {
	var op TemplateOp
	op.FS = testFS

	equal := func(expectedFile string, data []byte) bool {
		expectedData, err := op.readFile(expectedFile)
		if err != nil {
			t.Fatalf("%v", err)
		}
		n := len(data)
		if n != len(expectedData) {
			return false
		}
		for i, b := range data {
			if b != expectedData[i] {
				return false
			}
		}
		return true
	}
	verify := func(op *TemplateOp, templateFile, expectedFile string) {
		inputData, err := op.readFile(templateFile)
		if err != nil {
			t.Fatalf("%v", err)
		}
		values, err := op.Values()
		if err != nil {
			t.Fatalf("values: %v", err)
		}
		var buf bytes.Buffer
		tpl, err := op.buildTemplate(inputData)
		if err != nil {
			t.Fatalf("template %s: %v", templateFile, err)
		}
		err = tpl.Execute(&buf, values)
		if err != nil {
			t.Fatalf("execute: %v", err)
		}
		data := buf.Bytes()
		if !equal(expectedFile, data) {
			t.Fatalf("output (expected %s):\n%s", expectedFile, data)
		}
	}
	op.YamlFiles = []string{"test/properties.yaml"}
	verify(&op, "test/a.tpl", "test/a1.txt")

	op.YamlFiles = []string{"test/properties.yaml"}
	op.KeyValues = []string{"b=B2"}
	verify(&op, "test/a.tpl", "test/a2.txt")

	op.YamlFiles = []string{"test/properties.yaml"}
	op.KeyFiles = []string{"b=test/b.txt"}
	verify(&op, "test/a.tpl", "test/a3.txt")

	op.YamlFiles = []string{"test/properties.yaml"}
	op.JsonFiles = []string{"b=test/b.json"}
	verify(&op, "test/a.tpl", "test/a3.txt")

	op.JsonFiles = []string{"list=test/list.json"}
	verify(&op, "test/list.tpl", "test/list.txt")
}
