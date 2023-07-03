package stub

type StubProvider struct {
	Err error
}

func (tp *StubProvider) SendExchangeRate(rate float32, subscribers []string) error {
	return tp.Err
}
