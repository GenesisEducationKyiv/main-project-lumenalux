package rate

import "net/http"

type MockHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

type MockProvider struct {
	Rate  float32
	Error error
}

func (m *MockProvider) ExchangeRate() (float32, error) {
	return m.Rate, m.Error
}
