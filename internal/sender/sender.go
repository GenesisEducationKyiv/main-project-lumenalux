package sender

import "gses2-app/pkg/types"

type SenderProvider interface {
	SendExchangeRate(rate types.Rate, subscribers []types.Subscriber) error
}

type Service struct {
	provider SenderProvider
}

func NewService(provider SenderProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) SendExchangeRate(
	rate types.Rate,
	subscribers ...types.Subscriber,
) error {
	return s.provider.SendExchangeRate(rate, subscribers)
}
