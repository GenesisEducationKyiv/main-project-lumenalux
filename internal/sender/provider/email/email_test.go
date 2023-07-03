package email

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/pkg/config"
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
		dialer         smtp.TLSConnectionDialer
		factory        smtp.SMTPClientFactory
		expectedErr    error
	}{
		{
			name:           "Successful SendExchangeRate",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &smtp.StubDialer{},
			factory:        &smtp.StubSMTPClientFactory{Client: &smtp.StubSMTPClient{}},
			expectedErr:    nil,
		},
		{
			name:           "Failed due to dialer error",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &smtp.StubDialer{Err: errDialerError},
			factory:        &smtp.StubSMTPClientFactory{Client: &smtp.StubSMTPClient{}},
			expectedErr:    errDialerError,
		},
		{
			name:           "Failed due to factory error",
			emailAddresses: []string{"test@example.com"},
			exchangeRate:   10.5,
			dialer:         &smtp.StubDialer{},
			factory: &smtp.StubSMTPClientFactory{
				Client: &smtp.StubSMTPClient{},
				Err:    errFactoryError,
			},
			expectedErr: errFactoryError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			config := &config.Config{}
			service, err := NewProvider(config, tt.dialer, tt.factory)

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
