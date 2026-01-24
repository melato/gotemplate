package yaml

import (
	"gopkg.in/yaml.v2"
)

// A property parser that uses gopkg.in/yaml.v2
func ParseYaml(data []byte, properties map[string]any) error {
	return yaml.Unmarshal(data, &properties)
}
