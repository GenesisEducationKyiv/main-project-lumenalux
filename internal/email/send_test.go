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
			expectQuitCalled: true,
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
			expectQuitCalled: false,
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
			expectQuitCalled: false,
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

			if err != nil && tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
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
