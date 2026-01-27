package build

type Config struct {
	InputDir        string         `yaml:"input_dir"`
	OutputDir       string         `yaml:"output_dir"`
	InputExtension  string         `yaml:"input_ext"`
	OutputExtension string         `yaml:"output_ext"`
	Template        TemplateConfig `yaml:"template"`
	Properties      map[string]any `yaml:"properties"`
}

type TemplateConfig struct {
	Dir      string   `yaml:"dir,omitempty"`
	Patterns []string `yaml:"patterns,omitempty"`
}
