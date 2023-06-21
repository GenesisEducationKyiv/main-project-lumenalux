package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func Load(filename string) (Config, error) {
	var configuration Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return configuration, err
	}

	err = yaml.Unmarshal(data, &configuration)
	return configuration, err
}
