package gotemplate

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Options struct {
	YamlFiles []string `name:"f" usage:"yaml file"`
	JsonFiles []string `name:"json" usage:"key={json-file} - set the value of <key> to the content of <file>, parsed as JSON"`
	KeyValues []string `name:"D" usage:"key=value - set a property"`
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
	return builder.Unmarshal(key, data, unmarshal)
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

func (t *Options) addKeyValues(builder *builder, args []string) error {
	pairs, err := parseKeyValues(args, false)
	if err != nil {
		return err
	}
	for _, pair := range pairs {
		err := builder.Set(pair.Key, pair.Value)
		if err != nil {
			return err
		}
	}
	return nil
}

type keyValue struct {
	Key   string
	Value string
}

func parseKeyValues(args []string, allowMissingKeys bool) ([]keyValue, error) {
	pairs := make([]keyValue, len(args))
	for i, arg := range args {
		key, file, hasEq := strings.Cut(arg, "=")
		if hasEq {
			pairs[i] = keyValue{key, file}
		} else {
			if allowMissingKeys {
				pairs[i] = keyValue{"", arg}
			} else {
				return nil, fmt.Errorf("expected key=file: %s", arg)
			}
		}
	}
	return pairs, nil
}

func (t *Options) addEncodedFiles(builder *builder,
	unmarshal func([]byte, any) error,
	args []string) error {
	pairs, err := parseKeyValues(args, true)
	if err != nil {
		return err
	}
	for _, p := range pairs {
		data, err := t.readFile(p.Value)
		if err != nil {
			return err
		}
		return builder.Unmarshal(p.Key, data, unmarshal)
	}
	return nil
}

func (t *Options) apply(builder *builder) error {
	var err error
	if err == nil {
		err = t.addEncodedFiles(builder, yaml.Unmarshal, t.YamlFiles)
	}
	if err == nil {
		err = t.addEncodedFiles(builder, json.Unmarshal, t.JsonFiles)
	}
	if err == nil {
		err = t.addKeyValues(builder, t.KeyValues)
	}
	return err
}

func (t *Options) Values() (map[string]any, error) {
	var builder builder
	err := t.apply(&builder)
	if err != nil {
		return nil, err
	}
	return builder.values, nil
}
