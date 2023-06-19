package rate

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestKunaProviderExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		mockHTTPClient *MockHTTPClient
		expectedRate   float32
		expectedError  string
	}{
		{
			name: "Success",
			mockHTTPClient: &MockHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(
						bytes.NewBufferString(
							`[[123456789,"BTCUSDT","1.23","1.24","1.25","1.26",1.23,1.24,1.25]]`,
						),
					),
				},
				Error: nil,
			},
			expectedRate: 1.24,
		},
		{
			name: "HTTP request failure",
			mockHTTPClient: &MockHTTPClient{
				Response: nil,
				Error:    errors.New("http request failure"),
			},
			expectedError: "http request failure",
		},
		{
			name: "Unexpected status code",
			mockHTTPClient: &MockHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
				Error: nil,
			},
			expectedError: "unexpected status code",
		},
		{
			name: "Bad response body format",
			mockHTTPClient: &MockHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`[[]]`)),
				},
				Error: nil,
			},
			expectedError: "unexpected response format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := NewKunaProvider(tt.mockHTTPClient)
			rate, err := provider.ExchangeRate()

			if err != nil && tt.expectedError == "" {
				t.Errorf("didn't expect an error but got: %v", err)
			}

			if err != nil && !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("expected error message to contain %q, got %q", tt.expectedError, err.Error())
			}

			if err == nil && tt.expectedError != "" {
				t.Errorf("expected an error %q but got nil", tt.expectedError)
			}

			if rate != tt.expectedRate {
				t.Errorf("expected rate %v, got %v", tt.expectedRate, rate)
			}
		})
	}
}
