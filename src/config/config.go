package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	defaultCredentialsPath = "credentials/firebase.json"
	defaultStoragePath     = "storage.db"
	defaultSchemaPath      = "schema.json"
)

type Database struct {
	Storage string `yaml:"storage"`
	Schema  string `yaml:"schema"`
}

type Config struct {
	Url         string   `yaml:"url"`
	LogPath     string   `yaml:"log_path"`
	LogLevel    string   `yaml:"log_level"`
	Database    Database `yaml:"database"`
	Credentials string   `yaml:"credentials"`
}

// Sets default values if field is zero value.
func (config *Config) setDefaults() {
	if config.Database.Storage == "" {
		config.Database.Storage = defaultStoragePath
	}

	if config.Database.Schema == "" {
		config.Database.Schema = defaultSchemaPath
	}

	if config.Credentials == "" {
		config.Credentials = defaultCredentialsPath
	}
}

// Loads config from file.
func Load(path string) (*Config, error) {
	fmt.Printf("[Config] read config file: %s\n", path)

	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("[Config] [Error] failed read config file: %s\n", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("[Config] [Error] failed map config file: %s\n", err)
	}

	config.setDefaults()

	fmt.Println("[Config] config loaded successfully: ", *config)

	return config, nil
}
