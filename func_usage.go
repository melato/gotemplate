package gotemplate

import (
	"fmt"
	"strings"
)

type FuncUsage struct {
	Description string  `yaml:"desc"`
	Params      []Param `yaml:"params"`
}

func firstLine(s string) string {
	i := strings.IndexAny(s, "\r\n")
	if i >= 0 {
		return s[:i]
	} else {
		return s
	}
}

type Param struct {
	Name        string
	Description string
}

func (t Param) MarshalYAML() (interface{}, error) {
	return map[string]string{t.Name: t.Description}, nil
}

func (t *Param) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var m map[string]string
	err := unmarshal(&m)
	if err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("expected a single property: %v", m)
	}
	for key, value := range m {
		t.Name = key
		t.Description = value
	}
	return nil
}
