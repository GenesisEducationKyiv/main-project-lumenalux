package kuna

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/pkg/config"
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
		mockHTTPClient *StubHTTPClient
		expectedRate   float32
		expectedError  error
	}{
		{
			name: "Success",
			mockHTTPClient: &StubHTTPClient{
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
			mockHTTPClient: &StubHTTPClient{
				Response: nil,
				Error:    errors.New("http request failure"),
			},
			expectedError: ErrHTTPRequestFailure,
		},
		{
			name: "Unexpected status code",
			mockHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			expectedError: ErrUnexpectedStatusSode,
		},
		{
			name: "Bad response body format",
			mockHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[[]]`)),
				},
			},
			expectedError: ErrUnexpectedResponseFormat,
		},
		{
			name: "Bad response body format rate isn't a float64",
			mockHTTPClient: &StubHTTPClient{
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
			provider := NewKunaProvider(config, tt.mockHTTPClient)
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
