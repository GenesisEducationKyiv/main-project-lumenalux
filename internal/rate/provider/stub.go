package provider

type StubProvider struct {
	Rate  float32
	Error error
}

func (m *StubProvider) ExchangeRate() (float32, error) {
	return m.Rate, m.Error
}
