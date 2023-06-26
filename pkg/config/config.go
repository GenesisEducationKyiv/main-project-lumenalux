package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	ErrReadFile      = errors.New("failed to read config file")
	ErrUnmarshalYAML = errors.New("failed to unmarshal yaml data")
)

func Load(filename string) (Config, error) {
	var configuration Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return configuration, fmt.Errorf("%w: %s", ErrReadFile, err)
	}

	err = yaml.Unmarshal(data, &configuration)
	if err != nil {
		return configuration, fmt.Errorf("%w: %s", ErrUnmarshalYAML, err)
	}

	return configuration, nil
}
