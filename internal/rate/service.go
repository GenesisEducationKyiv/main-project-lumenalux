package rate

type Service interface {
	ExchangeRate() (float32, error)
}

type ServiceImpl struct {
	provider Provider
}

func NewService(provider Provider) Service {
	return &ServiceImpl{
		provider: provider,
	}
}

func (s *ServiceImpl) ExchangeRate() (float32, error) {
	return s.provider.ExchangeRate()
}
