package kuna

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/pkg/config"
	"gses2-app/pkg/types"
)

type StubHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *StubHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

func TestKunaProviderExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		stubHTTPClient *StubHTTPClient
		expectedRate   types.Rate
		expectedError  error
	}{
		{
			name: "Success",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(
							`[[123456789,"BTCUSDT","1.23","1.24","1.25","1.26",1.23,1.24,1.25]]`,
						),
					),
				},
			},
			expectedRate: 1.24,
		},
		{
			name: "HTTP request failure",
			stubHTTPClient: &StubHTTPClient{
				Response: nil,
				Error:    errors.New("http request failure"),
			},
			expectedError: ErrHTTPRequestFailure,
		},
		{
			name: "Unexpected status code",
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			expectedError: ErrUnexpectedStatusCode,
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
							`[[123456789,"BTCUSDT","1.23","1.24","1.25","1.26",1.23,true,1.25]]`,
						),
					),
				},
			},
			expectedError: ErrUnexpectedExchangeRateFormat,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := config.KunaAPIConfig{}
			logFunc := func(string, *http.Response) {}
			provider := NewProvider(config, tt.stubHTTPClient, logFunc)
			rate, err := provider.ExchangeRate()

			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedRate, rate, "Expected rate %v, got %v", tt.expectedRate, rate)
		})
	}

}
