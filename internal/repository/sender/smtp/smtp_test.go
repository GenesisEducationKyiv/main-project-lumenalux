package smtp

import (
	"errors"
	"testing"
)

var (
	errConnectionFailed = errors.New("failed to create connection")
	errSMTPClientFailed = errors.New("failed to create SMTP client")
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name        string
		dialer      TLSConnectionDialer
		factory     SMTPClientFactory
		config      SMTPConfig
		expectedErr error
	}{
		{
			name:    "Successful Connection",
			dialer:  &StubDialer{},
			factory: &StubSMTPClientFactory{Client: &StubSMTPClient{}},
			config: SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				User:     "user@example.com",
				Password: "password",
			},
			expectedErr: nil,
		},
		{
			name:    "Fail to create connection",
			dialer:  &StubDialer{Err: errConnectionFailed},
			factory: &StubSMTPClientFactory{Client: &StubSMTPClient{}},
			config: SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				User:     "user@example.com",
				Password: "password",
			},
			expectedErr: errConnectionFailed,
		},
		{
			name:   "Fail to create SMTP client",
			dialer: &StubDialer{},
			factory: &StubSMTPClientFactory{
				Client: &StubSMTPClient{}, Err: errSMTPClientFailed,
			},
			config: SMTPConfig{
				Host:     "smtp.example.com",
				Port:     587,
				User:     "user@example.com",
				Password: "password",
			},
			expectedErr: errSMTPClientFailed,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			client := NewSMTPClient(tt.config, tt.dialer, tt.factory)
			smtpClient, err := client.Connect()

			if err != nil && !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Connect() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			if err == nil && smtpClient == nil {
				t.Errorf("Expected smtp.Client, got nil")
			}
		})
	}
}
