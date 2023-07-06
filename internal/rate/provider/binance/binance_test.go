package binance

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

func TestBinanceProviderExchangeRate(t *testing.T) {
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
							`[[168,"153.00","153.0","153.0","123.456","0.0",99,"0.0",0,"0.0","0.0","0"]]`,
						),
					),
				},
			},
			expectedRate: 123.456,
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
			expectedError: ErrUnexpectedStatusSode,
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
							`[[168,"153.00","153.0","153.0",true,"0.0",99,"0.0",0,"0.0","0.0","0"]]`,
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

			config := config.BinanceAPIConfig{}
			logFunc := func(string, *http.Response) {}
			provider := NewProvider(config, tt.stubHTTPClient, logFunc)
			rate, err := provider.ExchangeRate()

			if tt.expectedError != nil {
				require.Error(t, err, "Expected an error but got nil")
				require.Contains(
					t, err.Error(), tt.expectedError.Error(),
					"Expected error message to contain %v, got %v", tt.expectedError, err,
				)
			} else {
				require.NoError(t, err, "Didn't expect an error but got: %v", err)
			}

			require.Equal(t, tt.expectedRate, rate, "Expected rate %v, got %v", tt.expectedRate, rate)
		})
	}

}
