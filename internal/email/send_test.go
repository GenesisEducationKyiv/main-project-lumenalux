package email

import (
	"errors"
	"testing"
)

type testCase struct {
	name             string
	client           *StubSenderSMTPClient
	email            *EmailMessage
	expectedErr      error
	expectDataCalled bool
	expectQuitCalled bool
}

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
			expectQuitCalled: true,
		},
		{
			name: "Error on write",
			client: &StubSenderSMTPClient{
				writeShouldReturn: errors.New("write error"),
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errors.New("write error"),
			expectDataCalled: true,
			expectQuitCalled: false,
		},
		{
			name: "Error on setMail",
			client: &StubSenderSMTPClient{
				mailShouldReturn: errors.New("set mail error"),
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errors.New("set mail error"),
			expectDataCalled: false,
			expectQuitCalled: false,
		},
		{
			name: "Error on setRecipients",
			client: &StubSenderSMTPClient{
				rcptShouldReturn: errors.New("set recipients error"),
			},
			email: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expectedErr:      errors.New("set recipients error"),
			expectDataCalled: false,
			expectQuitCalled: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SendEmail(tt.client, tt.email)
			if err != nil && tt.expectedErr == nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if err == nil && tt.expectedErr != nil {
				t.Error("Expected error, got nil")
			}

			if err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error() {
				t.Errorf("Error: got %v, want %v", err, tt.expectedErr)
			}

			if tt.client.dataCalled != tt.expectDataCalled {
				t.Errorf("Data called: got %v, want %v", tt.client.dataCalled, tt.expectDataCalled)
			}

			if tt.client.quitCalled != tt.expectQuitCalled {
				t.Errorf("Quit called: got %v, want %v", tt.client.quitCalled, tt.expectQuitCalled)
			}
		})
	}
}
