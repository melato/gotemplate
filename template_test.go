package gotemplate

import (
	"bytes"
	"embed"
	"testing"
)

//go:embed test/*
var testFS embed.FS

func verifyTemplate(t *testing.T, templateFile, expectedFile string, options Options) {
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
	op := TemplateOp{Options: options}
	op.DefineFuncs()
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

func TestOptions(t *testing.T) {
	verifyTemplate(t, "test/a.tpl", "test/a1.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})

	verifyTemplate(t, "test/a.tpl", "test/a1.txt", Options{
		JsonFiles: []string{"test/properties.json"},
	})
	verifyTemplate(t, "test/a.tpl", "test/a2.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		KeyValues: []string{"b=B2"},
	})

	verifyTemplate(t, "test/a.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		KeyFiles:  []string{"b=test/b.txt"},
	})
	verifyTemplate(t, "test/a.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
		JsonFiles: []string{"b=test/b.json"},
	})

	verifyTemplate(t, "test/list.tpl", "test/list.txt", Options{
		JsonFiles: []string{"list=test/list.json"},
	})
}

func TestFileFuncs(t *testing.T) {
	verifyTemplate(t, "test/func_file.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})
	verifyTemplate(t, "test/func_json.tpl", "test/a3.txt", Options{
		YamlFiles: []string{"test/properties.yaml"},
	})
	verifyTemplate(t, "test/func_yaml.tpl", "test/a1.txt", Options{})

}
