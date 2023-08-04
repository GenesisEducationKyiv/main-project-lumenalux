package sender

import (
	"gses2-app/internal/core/port"
)

type SenderPort interface {
	SendExchangeRate(rate port.Rate, subscribers []port.User) error
}

type Service struct {
	senderPort SenderPort
}

func NewService(provider SenderPort) *Service {
	return &Service{senderPort: provider}
}

func (s *Service) SendExchangeRate(
	rate port.Rate,
	users ...port.User,
) error {
	return s.senderPort.SendExchangeRate(rate, users)
}
