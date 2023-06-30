package rate

import (
	"errors"
	"testing"
)

func TestExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		mockProvider   *StubProvider
		expectedRate   float32
		expectingError bool
	}{
		{
			name: "Success",
			mockProvider: &StubProvider{
				Rate:  1.23,
				Error: nil,
			},
			expectedRate:   1.23,
			expectingError: false,
		},
		{
			name: "Failure",
			mockProvider: &StubProvider{
				Rate:  0,
				Error: errors.New("error fetching rate"),
			},
			expectedRate:   0,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.mockProvider)
			rate, err := service.ExchangeRate()

			if rate != tt.expectedRate {
				t.Errorf("expected rate %v, got %v", tt.expectedRate, rate)
			}

			if tt.expectingError && err == nil {
				t.Errorf("expected an error but got nil")
			}

			if !tt.expectingError && err != nil {
				t.Errorf("didn't expect an error but got: %v", err)
			}
		})
	}
}
