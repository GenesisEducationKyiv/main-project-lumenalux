package stub

import "gses2-app/pkg/types"

type StubProvider struct {
	Err error
}

func (tp *StubProvider) SendExchangeRate(rate types.Rate, subscribers []types.User) error {
	return tp.Err
}
