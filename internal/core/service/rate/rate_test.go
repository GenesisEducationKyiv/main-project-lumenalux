package rate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
)

type StubLogger struct{}

func (s *StubLogger) Info(...interface{})           {}
func (s *StubLogger) Infof(string, ...interface{})  {}
func (s *StubLogger) Debug(...interface{})          {}
func (s *StubLogger) Debugf(string, ...interface{}) {}
func (s *StubLogger) Error(...interface{})          {}
func (s *StubLogger) Errorf(string, ...interface{}) {}

type StubProvider struct {
	Rate         port.Rate
	Error        error
	ProviderName string
}

func (m *StubProvider) ExchangeRate() (port.Rate, error) {
	return m.Rate, m.Error
}

func (m *StubProvider) Name() string {
	return m.ProviderName
}

func TestExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		stubProvider   *StubProvider
		expectedRate   port.Rate
		expectingError bool
	}{
		{
			name: "Success",
			stubProvider: &StubProvider{
				Rate:  1.23,
				Error: nil,
			},
			expectedRate:   1.23,
			expectingError: false,
		},
		{
			name: "Failure",
			stubProvider: &StubProvider{
				Rate:  0,
				Error: errors.New("error fetching rate"),
			},
			expectedRate:   0,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			service := NewService(&StubLogger{}, tt.stubProvider)
			rate, err := service.ExchangeRate()

			require.Equal(
				t, tt.expectedRate, rate,
				"Expected rate %v, got %v", tt.expectedRate, rate,
			)

			require.Equal(t, tt.expectingError, err != nil)
		})
	}

}
