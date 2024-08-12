package config

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	ini "gopkg.in/ini.v1"
	yaml "gopkg.in/yaml.v2"
)

// ConfigFileReader represent a reader to read a configuration file.
// It also store the root path for the configuration files.
type ConfigFileReader struct {
	root string
}

// NewConfigFileReader returns an instance of a ConfigFileReader
// with root as a root of file location/directory
func NewConfigFileReader(root string) *ConfigFileReader {
	return &ConfigFileReader{
		root: root,
	}
}

// Read reads a configuration file from the given root+filePath location.
// It returns config.ErrNoFileFound if the given root+filePath is not exist.
// It also satisfy FileReader interface for method FileReader.Read
func (cfr *ConfigFileReader) Read(ctx context.Context, dest interface{}, filePath string) error {
	log.Println(ctx, nil, nil, "reading from config file: %s - ConfigFileReader.Read", filePath)
	return read(dest, cfr.root+filePath)
}

// Read configuration from the given paths.
// it will use the first file found in the given paths.
// It returns ErrNoFileFound if none of the given paths is exist
func read(dest interface{}, paths ...string) error {
	for _, path := range paths {
		path = replacePathByEnv(path)

		// check if this path is exist
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		// load config
		ext := filepath.Ext(path)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		switch {
		case ext == ".ini":
			return loadIniConfig(dest, path)
		case ext == ".yaml" || ext == ".yml":
			return yaml.Unmarshal(content, dest)
		case ext == ".json":
			return json.Unmarshal(content, dest)
		}
	}
	return errors.New("file not found")
}

func loadIniConfig(dest interface{}, path string) error {
	f, err := ini.Load(path)
	if err != nil {
		return err
	}
	return f.MapTo(dest)
}

func replacePathByEnv(path string) string {
	return path
}
