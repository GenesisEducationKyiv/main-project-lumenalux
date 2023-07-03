package sender

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	name             string
	client           *StubSenderSMTPClient
	email            *EmailMessage
	expectedErr      error
	expectDataCalled bool
}

var (
	errWrite         = errors.New("write error")
	errSetMail       = errors.New("set mail error")
	errSetRecipients = errors.New("set recipients error")
)

func TestSendEmail(t *testing.T) {
	tests := []testCase{
		{
			name: "Send email",
			client: &StubSenderSMTPClient{
				writeShouldReturn: nil,
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      nil,
			expectDataCalled: true,
		},
		{
			name: "Error on write",
			client: &StubSenderSMTPClient{
				writeShouldReturn: errWrite,
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errWrite,
			expectDataCalled: true,
		},
		{
			name: "Error on setMail",
			client: &StubSenderSMTPClient{
				mailShouldReturn: errSetMail,
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errSetMail,
			expectDataCalled: false,
		},
		{
			name: "Error on setRecipients",
			client: &StubSenderSMTPClient{
				rcptShouldReturn: errSetRecipients,
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errSetRecipients,
			expectDataCalled: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := SendEmail(tt.client, tt.email)

			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr, "Error: got %v, want %v", err, tt.expectedErr)
			} else {
				require.NoError(t, err, "Unexpected error: %v", err)
			}

			require.Equal(t, tt.expectDataCalled, tt.client.dataCalled, "Data called: got %v, want %v", tt.client.dataCalled, tt.expectDataCalled)
		})
	}
}
