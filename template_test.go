package gotemplate

import (
	"bytes"
	"embed"
	"testing"
)

//go:embed test/*
var testFS embed.FS

func TestOptions(t *testing.T) {
	equal := func(expectedData, data []byte) bool {
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
	verify := func(templateFile, expectedFile string, options Options) {
		op := TemplateOp{Options: options}
		op.Init()
		op.FS = testFS
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
		expectedData, err := op.readFile(expectedFile)
		if err != nil {
			t.Fatalf("%v", err)
		}
		data := buf.Bytes()
		if !equal(expectedData, data) {
			t.Fatalf("expected %s (%d bytes), got %d bytes:\n%s",
				expectedFile,
				len(expectedData),
				len(data),
				data)
		}
	}
	verify("test/a.tpl", "test/a1.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})

	verify("test/a.tpl", "test/a1.txt", Options{
		JsonFiles: []string{"test/properties.json"},
	})
	verify("test/a.tpl", "test/a2.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		KeyValues: []string{"b=B2"},
	})

	verify("test/a.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		KeyFiles:  []string{"b=test/b.txt"},
	})
	verify("test/func_file.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})
	verify("test/func_json.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})
	verify("test/func_yaml.tpl", "test/a1.txt", Options{})

	verify("test/a.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		JsonFiles: []string{"b=test/b.json"},
	})

	verify("test/list.tpl", "test/list.txt", Options{
		JsonFiles: []string{"list=test/list.json"},
	})
}
