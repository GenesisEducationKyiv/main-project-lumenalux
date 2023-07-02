package sender

import (
	"errors"
	"testing"

	"gses2-app/pkg/config"

	"github.com/stretchr/testify/require"
)

var (
	errDialerError  = errors.New("dialer error")
	errFactoryError = errors.New("factory error")
)

func TestSendExchangeRate(t *testing.T) {
	tests := []struct {
		name           string
		emailAddresses []string
		exchangeRate   float32
		dialer         TLSConnectionDialer
		factory        SMTPClientFactory
		expectedErr    error
	}{
		{
			name:           "Successful SendExchangeRate",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &StubDialer{},
			factory:        &StubSMTPClientFactory{Client: &StubSMTPClient{}},
			expectedErr:    nil,
		},
		{
			name:           "Failed due to dialer error",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &StubDialer{Err: errDialerError},
			factory:        &StubSMTPClientFactory{Client: &StubSMTPClient{}},
			expectedErr:    errDialerError,
		},
		{
			name:           "Failed due to factory error",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &StubDialer{},
			factory: &StubSMTPClientFactory{
				Client: &StubSMTPClient{},
				Err:    errFactoryError,
			},
			expectedErr: errFactoryError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &config.Config{}
			service, err := NewSenderService(config, tt.dialer, tt.factory)

			if tt.expectedErr != nil {
				require.ErrorIs(
					t,
					err,
					tt.expectedErr,
					"NewSenderService() error = %v, expectedErr %v",
					err,
					tt.expectedErr,
				)
			} else {
				require.NoError(t, err, "NewSenderService() unexpected error = %v", err)
			}

			if tt.expectedErr != nil {
				return
			}

			err = service.SendExchangeRate(tt.exchangeRate, tt.emailAddresses)

			if tt.expectedErr != nil {
				require.ErrorIs(
					t,
					err,
					tt.expectedErr,
					"SendExchangeRate() error = %v, expectedErr %v",
					err,
					tt.expectedErr,
				)
			} else {
				require.NoError(t, err, "SendExchangeRate() unexpected error = %v", err)
			}
		})
	}
}
