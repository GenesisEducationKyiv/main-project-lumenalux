package stub

import "gses2-app/pkg/types"

type StubProvider struct {
	Rate  types.Rate
	Error error
}

func (m *StubProvider) ExchangeRate() (types.Rate, error) {
	return m.Rate, m.Error
}

func (m *StubProvider) Name() string {
	return "StubRateProvider"
}
