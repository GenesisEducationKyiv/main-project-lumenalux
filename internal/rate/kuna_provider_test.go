package rate

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"gses2-app/pkg/config"
)

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
		t.Run(tt.name, func(t *testing.T) {
			config := config.KunaAPIConfig{}
			provider := NewKunaProvider(config, tt.mockHTTPClient)
			rate, err := provider.ExchangeRate()

			if err != nil && tt.expectedError == nil {
				t.Errorf("didn't expect an error but got: %v", err)
			}

			if err != nil && !strings.Contains(err.Error(), tt.expectedError.Error()) {
				t.Errorf("expected error message to contain %v, got %v", tt.expectedError, err)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected an error %v but got nil", tt.expectedError)
			}

			if rate != tt.expectedRate {
				t.Errorf("expected rate %v, got %v", tt.expectedRate, rate)
			}
		})
	}
}
