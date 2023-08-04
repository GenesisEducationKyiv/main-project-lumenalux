package coingecko

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/rate/rest"
)

type StubLogger struct{}

func (s *StubLogger) Info(...interface{})           {}
func (s *StubLogger) Infof(string, ...interface{})  {}
func (s *StubLogger) Debug(...interface{})          {}
func (s *StubLogger) Debugf(string, ...interface{}) {}
func (s *StubLogger) Error(...interface{})          {}
func (s *StubLogger) Errorf(string, ...interface{}) {}

type StubHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *StubHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

func TestBinanceProviderExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		stubHTTPClient *StubHTTPClient
		expectedRate   port.Rate
		expectedError  error
	}{
		{
			name: "Success",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(
							`{"bitcoin":{"uah":123456}}`,
						),
					),
				},
			},
			expectedRate: 123456,
		},
		{
			name: "HTTP request failure",
			stubHTTPClient: &StubHTTPClient{
				Response: nil,
				Error:    rest.ErrHTTPRequestFailure,
			},
			expectedError: rest.ErrHTTPRequestFailure,
		},
		{
			name: "Unexpected status code",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			expectedError: rest.ErrUnexpectedStatusCode,
		},
		{
			name: "Bad response body format",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[[]]`)),
				},
			},
			expectedError: ErrUnexpectedResponseFormat,
		},
		{
			name: "Bad response body format rate isn't a float64",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(
							`{"bitcoin":{"uah":false}}`,
						),
					),
				},
			},
			expectedError: ErrUnexpectedResponseFormat,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := CoingeckoAPIConfig{}
			provider := NewProvider(&StubLogger{}, config, tt.stubHTTPClient)
			rate, err := provider.ExchangeRate()

			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedRate, rate, "Expected rate %v, got %v", tt.expectedRate, rate)
		})
	}
}
