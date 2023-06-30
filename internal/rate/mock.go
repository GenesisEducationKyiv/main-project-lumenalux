package rate

import "net/http"

type StubHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *StubHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

type StubProvider struct {
	Rate  float32
	Error error
}

func (m *StubProvider) ExchangeRate() (float32, error) {
	return m.Rate, m.Error
}
