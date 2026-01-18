package funcs

import (
	"os"
)

// ReadFile reads a file from the filesystem, and returns it as a string
func ReadFile(file string) (string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
