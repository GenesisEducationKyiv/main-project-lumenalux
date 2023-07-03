package sender

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/sender/provider/stub"
)

var (
	errProvider = errors.New("provider error")
)

func TestService_SendExchangeRate(t *testing.T) {
	tests := []struct {
		name        string
		providerErr error
		expectedErr error
	}{
		{
			name:        "No error from provider",
			providerErr: nil,
			expectedErr: nil,
		},
		{
			name:        "Error from provider",
			providerErr: errProvider,
			expectedErr: errProvider,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			provider := &stub.StubProvider{Err: tt.providerErr}
			service := NewService(provider)

			err := service.SendExchangeRate(1.23, "subscriber")

			require.Equal(t, tt.expectedErr, err)
		})
	}
}
