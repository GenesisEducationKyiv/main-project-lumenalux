package provider

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"gses2-app/pkg/types"

	"github.com/stretchr/testify/require"
)

type StubHTTPClient struct {
	Response *http.Response
	Error    error
}

func (m *StubHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Error
}

type StubProvider struct {
	Url          string
	ProviderName string
	Rate         types.Rate
	Error        error
}

func (s *StubProvider) URL() string {
	return s.Url
}

func (s *StubProvider) Name() string {
	return s.ProviderName
}

func (s *StubProvider) ExtractRate(r *http.Response) (types.Rate, error) {
	return s.Rate, s.Error
}

func TestExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		stubProvider   Provider
		stubHTTPClient *StubHTTPClient
		expectedRate   types.Rate
		expectedError  error
	}{
		{
			name: "Success",
			stubProvider: &StubProvider{
				Url:          "https://test.url",
				ProviderName: "Test",
				Rate:         1.23,
			},
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString("Success Response")),
				},
			},
			expectedRate: 1.23,
		},
		{
			name: "HTTP request failure",
			stubProvider: &StubProvider{
				Url: "https://test.url",
			},
			stubHTTPClient: &StubHTTPClient{
				Error: ErrHTTPRequestFailure,
			},
			expectedError: ErrHTTPRequestFailure,
		},
		{
			name: "Unexpected status code",
			stubProvider: &StubProvider{
				Url: "https://test.url",
			},
			stubHTTPClient: &StubHTTPClient{
				Response: &http.Response{
					StatusCode: http.StatusForbidden,
				},
			},
			expectedError: ErrUnexpectedStatusCode,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			logFunc := func(providerName string, resp *http.Response) {}
			abstractProvider := NewProvider(tt.stubProvider, tt.stubHTTPClient, logFunc)
			rate, err := abstractProvider.ExchangeRate()

			require.ErrorIs(t, err, tt.expectedError)
			require.Equal(t, tt.expectedRate, rate, "Expected rate %v, got %v", tt.expectedRate, rate)
		})
	}
}
