package rate

import (
	"gses2-app/pkg/types"
	"log"
)

type Provider interface {
	ExchangeRate() (types.Rate, error)
	Name() string
}

type Service struct {
	providers []Provider
}

func NewService(providers ...Provider) *Service {
	return &Service{
		providers: providers,
	}
}

func (s *Service) ExchangeRate() (rate types.Rate, err error) {
	for _, provider := range s.providers {
		rate, err = provider.ExchangeRate()
		if err == nil {
			return rate, nil
		}

		log.Printf("Error, %v: %v", provider.Name(), err)
	}

	return rate, err
}
