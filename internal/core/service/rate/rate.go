package rate

import (
	"log"

	"gses2-app/internal/core/port"
)

type RatePort interface {
	ExchangeRate() (port.Rate, error)
	Name() string
}

type Service struct {
	providers []RatePort
}

func NewService(providers ...RatePort) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) ExchangeRate() (rate port.Rate, err error) {
	for _, provider := range s.providers {
		rate, err = provider.ExchangeRate()
		if err == nil {
			return rate, nil
		}

		log.Printf("Error, %v: %v", provider.Name(), err)
	}

	return rate, err
}
