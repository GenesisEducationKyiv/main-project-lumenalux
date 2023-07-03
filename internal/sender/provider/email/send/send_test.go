package send

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/sender/provider/email/message"
)

type testCase struct {
	name             string
	client           *StubSMTPClient
	email            *message.EmailMessage
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
			client: &StubSMTPClient{
				writeShouldReturn: nil,
			},
			email: &message.EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expectedErr:      nil,
			expectDataCalled: true,
		},
		{
			name: "Error on write",
			client: &StubSMTPClient{
				writeShouldReturn: errWrite,
			},
			email: &message.EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expectedErr:      errWrite,
			expectDataCalled: true,
		},
		{
			name: "Error on setMail",
			client: &StubSMTPClient{
				mailShouldReturn: errSetMail,
			},
			email: &message.EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expectedErr:      errSetMail,
			expectDataCalled: false,
		},
		{
			name: "Error on setRecipients",
			client: &StubSMTPClient{
				rcptShouldReturn: errSetRecipients,
			},
			email: &message.EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
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