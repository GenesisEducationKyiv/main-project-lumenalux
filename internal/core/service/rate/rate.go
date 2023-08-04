package rate

import (
	"gses2-app/internal/core/port"
)

type RatePort interface {
	ExchangeRate() (port.Rate, error)
	Name() string
}

type Service struct {
	providers []RatePort
	logger    port.Logger
}

func NewService(logger port.Logger, providers ...RatePort) *Service {
	return &Service{
		logger:    logger,
		providers: providers,
	}
}

func (s *Service) ExchangeRate() (rate port.Rate, err error) {
	for _, provider := range s.providers {
		rate, err = provider.ExchangeRate()
		if err == nil {
			return rate, nil
		}

		s.logger.Errorf("Error, %v: %v", provider.Name(), err)
	}

	return rate, err
}
