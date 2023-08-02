package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

var (
	ErrLoadEnvVariable = errors.New("failed to load env variables")
)

func Load(prefix string) (configuration Config, err error) {
	if err := envconfig.Process(prefix, &configuration); err != nil {
		return configuration, ErrLoadEnvVariable
	}

	return configuration, nil
}
