package config

import (
	"context"
	"errors"

	"github.com/sethvargo/go-envconfig"
)

var (
	ErrLoadEnvVariable = errors.New("failed to load env variables")
)

func Load(ctx context.Context) (configuration Config, err error) {
	if err := envconfig.Process(ctx, &configuration); err != nil {
		return configuration, ErrLoadEnvVariable
	}

	return configuration, nil
}
