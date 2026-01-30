package gotemplate

import (
	"gopkg.in/yaml.v2"
)

// A property parser that uses gopkg.in/yaml.v2
func ParseYaml(data []byte) (map[string]any, error) {
	var m map[string]any
	err := yaml.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
