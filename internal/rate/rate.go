package rate

import "gses2-app/pkg/types"

type Provider interface {
	ExchangeRate() (types.Rate, error)
}

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) ExchangeRate() (types.Rate, error) {
	return s.provider.ExchangeRate()
}
