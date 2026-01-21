package gotemplate

type builder struct {
	values map[string]any
}

func (t *builder) Set(key string, value any) {
	if t.values == nil {
		t.values = make(map[string]any)
	}
	t.values[key] = value
}

func (t *builder) Unmarshal(data []byte,
	unmarshal func([]byte, any) error) error {
	if t.values == nil {
		t.values = make(map[string]any)
	}
	return unmarshal(data, &t.values)
}
