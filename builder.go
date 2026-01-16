package gotemplate

import (
	"fmt"
)

type builder struct {
	values map[string]any
	keys   map[string]bool
}

func (t *builder) Set(key string, value any) error {
	if t.keys == nil {
		t.keys = make(map[string]bool)
	}
	_, exists := t.keys[key]
	t.keys[key] = true
	if exists {
		if key != "" {
			return fmt.Errorf("key %s is used more than once", key)
		} else {
			return fmt.Errorf("values are replaced more than once")
		}
	}
	if t.values == nil {
		t.values = make(map[string]any)
	}
	if key != "" {
		t.values[key] = value
	} else {
		m, isMap := value.(map[string]any)
		if !isMap {
			return fmt.Errorf("not a map: %T", value)
		}
		for key, value := range m {
			t.values[key] = value
		}
	}
	return nil
}

func (t *builder) Unmarshal(key string, data []byte,
	unmarshal func([]byte, any) error) error {
	if key == "" {
		t.values = make(map[string]any)
		return unmarshal(data, &t.values)
	} else {
		var v any
		err := unmarshal(data, &v)
		if err != nil {
			return err
		}
		return t.Set(key, v)
	}

}
