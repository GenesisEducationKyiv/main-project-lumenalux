package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func Load(filename string) (Config, error) {
	var configuration Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return configuration, fmt.Errorf("failed to read config file %s: %w", filename, err)
	}

	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return configuration, fmt.Errorf("failed to unmarshal yaml data: %w", err)
	}

	return configuration, nil
}
