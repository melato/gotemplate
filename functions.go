package gotemplate

import (
	"encoding/json"
	"os"

	"gopkg.in/yaml.v2"
)

type FileFunctions struct {
}

func (t *FileFunctions) File(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (t *FileFunctions) Json(file string) (any, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var v any
	err = json.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (t *FileFunctions) Yaml(file string) (map[string]any, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var v map[string]any
	err = yaml.Unmarshal(data, &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
