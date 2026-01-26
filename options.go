package gotemplate

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

type Options struct {
	PropertyFiles []string `name:"f" usage:"properties file"`
	Format        string   `name:"format" usage:"property file format"`
	KeyValues     []string `name:"D" usage:"key=value - set a property"`
	// if FS is not null, read files from FS, otherwise use os.ReadFile()
	// used for testing
	FS fs.FS
}

func (t *Options) readFile(file string) ([]byte, error) {
	if file == "" {
		return nil, fmt.Errorf("empty filename")
	}
	if t.FS == nil {
		return os.ReadFile(file)
	} else {
		return fs.ReadFile(t.FS, file)
	}
}

type unmarshalF func([]byte, any) error

func (t *Options) unmarshalFile(builder *builder, key string, file string, unmarshal unmarshalF) error {
	data, err := t.readFile(file)
	if err != nil {
		return err
	}
	return builder.Unmarshal(data, unmarshal)
}

// ParseKeyValue - parse a string of the form <key1[.key2]...>=<value>
// The keys are separated by dots
// The first returned argument is the keys, and the second is the value
// If there is no "=", then the keys are nil and the value is the input string
func (t *Options) ParseKeyValue(keyValue string) ([]string, string) {
	kv := strings.SplitN(keyValue, "=", 2)
	if len(kv) != 2 {
		return nil, keyValue
	}
	compoundKey, value := kv[0], kv[1]
	keys := strings.Split(compoundKey, ".")
	return keys, value
}

func (t *Options) addKeyValues(values map[any]any, args []string) error {
	pairs, err := parseKeyValues(args)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		values[pair.Key] = pair.Value
	}
	return nil
}

type keyValue struct {
	Key   string
	Value string
}

func parseKeyValues(args []string) ([]keyValue, error) {
	pairs := make([]keyValue, len(args))
	for i, arg := range args {
		key, file, hasEq := strings.Cut(arg, "=")
		if hasEq {
			pairs[i] = keyValue{key, file}
		} else {
			return nil, fmt.Errorf("expected key=file: %s", arg)
		}
	}
	return pairs, nil
}

func (t *Options) addEncodedFiles(values map[any]any,
	parse func([]byte, map[any]any) error,
	files []string) error {
	for _, file := range files {
		data, err := t.readFile(file)
		if err != nil {
			return err
		}
		err = parse(data, values)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Options) apply(values map[any]any) error {
	var err error
	parse := ParseProperties
	format := t.Format
	var foundParser bool
	if format != "" {
		parse, foundParser = propertyParsers[format]
		if !foundParser {
			return fmt.Errorf("unknown properties format: %s", format)
		}
	}
	if err == nil {
		err = t.addEncodedFiles(values, parse, t.PropertyFiles)
	}
	if err == nil {
		err = t.addKeyValues(values, t.KeyValues)
	}
	return err
}

func (t *Options) Values() (map[any]any, error) {
	values := make(map[any]any)
	err := t.apply(values)
	if err != nil {
		return nil, err
	}
	return values, nil
}
