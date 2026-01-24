package gotemplate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type parsers map[string]func([]byte, map[string]any) error

var propertyParsers parsers = make(parsers)

func SetParser(name string, parse func([]byte, map[string]any) error) {
	propertyParsers[name] = parse
}

func init() {
	SetParser("properties", ParseProperties)
	SetParser("json", ParseJson)
}

func ParseProperties(data []byte, properties map[string]any) error {
	for line := range iterLines(bytes.NewReader(data)) {
		line = strings.TrimLeft(line, " \t")
		if line == "" || line[0] == '#' {
			continue
		}
		name, value, found := strings.Cut(line, ":")
		if !found {
			return fmt.Errorf("cannot parse property: %s\n", line)
		}
		name = strings.TrimSpace(name)
		value = strings.TrimSpace(value)
		properties[name] = value
	}
	return nil
}

func ParseJson(data []byte, properties map[string]any) error {
	return json.Unmarshal(data, &properties)
}
