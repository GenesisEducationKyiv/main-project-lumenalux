package rate

type Provider interface {
	ExchangeRate() (float32, error)
}

type Service struct {
	provider Provider
}

func NewService(provider Provider) *Service {
	return &Service{
		provider: provider,
	}
}

func (s *Service) ExchangeRate() (float32, error) {
	return s.provider.ExchangeRate()
}
