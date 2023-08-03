package sender

import (
	"gses2-app/internal/rate"
	"gses2-app/internal/user/repository"
)

type SenderProvider interface {
	SendExchangeRate(rate rate.Rate, subscribers []repository.User) error
}

type Service struct {
	provider SenderProvider
}

func NewService(provider SenderProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) SendExchangeRate(
	rate rate.Rate,
	subscribers ...repository.User,
) error {
	return s.provider.SendExchangeRate(rate, subscribers)
}
