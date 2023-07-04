package rate

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/rate/provider/stub"
	"gses2-app/pkg/types"
)

func TestExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		stubProvider   *stub.StubProvider
		expectedRate   types.Rate
		expectingError bool
	}{
		{
			name: "Success",
			stubProvider: &stub.StubProvider{
				Rate:  1.23,
				Error: nil,
			},
			expectedRate:   1.23,
			expectingError: false,
		},
		{
			name: "Failure",
			stubProvider: &stub.StubProvider{
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
			service := NewService(tt.stubProvider)
			rate, err := service.ExchangeRate()

			require.Equal(
				t, tt.expectedRate, rate,
				"Expected rate %v, got %v", tt.expectedRate, rate,
			)

			if tt.expectingError {
				require.Error(t, err, "Expected an error but got nil")
				return
			}

			require.NoError(t, err, "Didn't expect an error but got: %v", err)
		})
	}

}
