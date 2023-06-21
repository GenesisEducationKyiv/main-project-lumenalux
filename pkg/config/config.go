package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var configuration Config

func Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &configuration)
}

func Current() Config {
	return configuration
}
