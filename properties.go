package gotemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type PropertyParser func([]byte) (map[string]any, error)
type parsers map[string]PropertyParser

var propertyParsers parsers = make(parsers)

func SetParser(name string, parse PropertyParser) {
	propertyParsers[name] = parse
}

func init() {
	SetParser("properties", ParseProperties)
	SetParser("json", ParseJson)
}

func ParseProperties(data []byte) (map[string]any, error) {
	properties := make(map[string]any)
	for line := range iterLines(bytes.NewReader(data)) {
		line = strings.TrimLeft(line, " \t")
		if line == "" || line[0] == '#' {
			continue
		}
		name, value, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("cannot parse property: %s\n", line)
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		properties[name] = value
	}
	return properties, nil
}

func ParseJson(data []byte) (map[string]any, error) {
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
