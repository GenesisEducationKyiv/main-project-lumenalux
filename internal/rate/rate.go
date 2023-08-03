package rate

import (
	"log"
)

// Rate represents the exchange rate between two currencies.
// It is expressed as a float32 value.
type Rate float32

type Provider interface {
	ExchangeRate() (Rate, error)
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

func (s *Service) ExchangeRate() (rate Rate, err error) {
	for _, provider := range s.providers {
		rate, err = provider.ExchangeRate()
		if err == nil {
			return rate, nil
		}

		log.Printf("Error, %v: %v", provider.Name(), err)
	}

	return rate, err
}
