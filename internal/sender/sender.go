package sender

type SenderProvider interface {
	SendExchangeRate(rate float32, subscribers []string) error
}

type Service struct {
	provider SenderProvider
}

func NewService(provider SenderProvider) *Service {
	return &Service{provider: provider}
}

func (s *Service) SendExchangeRate(
	rate float32,
	subscribers ...string,
) error {
	return s.provider.SendExchangeRate(rate, subscribers)
}
