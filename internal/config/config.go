package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Token    string   `yaml:"token"`
	URL      string   `yaml:"ip_url"`
	Services []string `yaml:"services,flow"`
}

func Load() (Config, error) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
