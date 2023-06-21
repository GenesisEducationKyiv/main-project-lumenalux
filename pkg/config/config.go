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

	err = yaml.Unmarshal(data, &configuration)
	return err
}

func Current() Config {
	return configuration
}
