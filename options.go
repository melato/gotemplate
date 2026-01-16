package gotemplate

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Options struct {
	YamlFiles []string `name:"f" usage:"yaml file"`
	JsonFiles []string `name:"json" usage:"key={json-file} - set the value of <key> to the content of <file>, parsed as JSON"`
	KeyFiles  []string `name:"F" usage:"key=file - set the value of <key> to the content of <file>"`
	KeyValues []string `name:"D" usage:"key=value - set a property"`
	FS        fs.FS
}

func (t *Options) readFile(file string) ([]byte, error) {
	if t.FS == nil {
		return os.ReadFile(file)
	} else {
		return fs.ReadFile(t.FS, file)
	}
}

func (t *Options) SetFile(properties Properties, keys []string, file string) error {
	data, err := t.readFile(file)
	if err != nil {
		return err
	}
	var fileProperties Properties
	err = yaml.Unmarshal(data, &fileProperties)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		for key, value := range fileProperties {
			properties[key] = value
		}
		return nil
	}
	return properties.Set(keys, fileProperties)
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

func (t *Options) Apply(properties Properties) error {
	for _, kvFile := range t.YamlFiles {
		keys, file := t.ParseKeyValue(kvFile)
		err := t.SetFile(properties, keys, file)
		if err != nil {
			return err
		}
	}
	for _, kv := range t.KeyValues {
		keys, value := t.ParseKeyValue(kv)
		if len(keys) == 0 {
			return fmt.Errorf("missing key(s): %s", kv)
		}
		err := properties.Set(keys, value)
		if err != nil {
			return err
		}
	}
	return nil

}

type KeyFile struct {
	Key  string
	File string
}

func parseKeyFiles(args []string) []KeyFile {
	result := make([]KeyFile, 0, len(args))
	for _, arg := range args {
		key, file, hasEq := strings.Cut(arg, "=")
		var k KeyFile
		if hasEq {
			k.Key = key
			k.File = file
		} else {
			k.Key = filepath.Base(arg)
			k.File = arg
		}
		result = append(result, k)
	}
	return result
}

func (t *Options) addKeyFiles(properties Properties, args []string) error {
	pairs := parseKeyFiles(args)
	for _, p := range pairs {
		data, err := os.ReadFile(p.File)
		if err != nil {
			return err
		}
		properties[p.Key] = string(data)
	}
	return nil
}

func (t *Options) addJsonFiles(properties Properties, args []string) error {
	pairs := parseKeyFiles(args)
	for _, p := range pairs {
		data, err := os.ReadFile(p.File)
		if err != nil {
			return err
		}
		var v any
		err = json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		properties[p.Key] = v
	}
	return nil
}

func (t *Options) Values() (Properties, error) {
	properties := make(Properties)
	err := t.Apply(properties)
	if err == nil {
		err = t.addKeyFiles(properties, t.KeyFiles)
	}
	if err == nil {
		err = t.addJsonFiles(properties, t.JsonFiles)
	}
	if err != nil {
		return nil, err
	}
	return properties, nil
}
