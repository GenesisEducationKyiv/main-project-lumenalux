package stub

import "gses2-app/pkg/types"

type StubProvider struct {
	Err error
}

func (tp *StubProvider) SendExchangeRate(rate types.Rate, subscribers []types.Subscriber) error {
	return tp.Err
}
