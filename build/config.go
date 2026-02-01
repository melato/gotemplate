package build

import (
	"path/filepath"
)

type Config struct {
	InputDir        string           `yaml:"input_dir"`
	OutputDir       string           `yaml:"output_dir"`
	InputExtension  string           `yaml:"input_ext"`
	OutputExtension string           `yaml:"output_ext"`
	Templates       []TemplateConfig `yaml:"templates"`
	Properties      map[string]any   `yaml:"properties"`
}

type TemplateConfig struct {
	Dir      string   `yaml:"dir,omitempty"`
	Patterns []string `yaml:"patterns,omitempty"`
}

// Resolve a path relative to a directory
func ResolvePath(dir string, path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(dir, path)

}
