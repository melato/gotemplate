package gotemplate

// Programmatic configuration of templates
// Use to add funcs to the FuncMap
type Config struct {
	BaseConfig
	funcUsage    map[string]string
	funcUsageTxt [][]byte
	parsedUsage  bool
}

/*
Add usage for functions
*/
func (t *Config) AddUsage(funcUsage map[string]string) {
	if t.funcUsage == nil {
		t.funcUsage = make(map[string]string)
	}
	for name, u := range funcUsage {
		t.funcUsage[name] = u
	}
}

/*
Add usage for functions, in text format
The usage is parsed when needed
*/
func (t *Config) AddUsageTxt(usageTxt []byte) {
	t.funcUsageTxt = append(t.funcUsageTxt, usageTxt)
}
