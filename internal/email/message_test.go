package email

import (
	"testing"

	"gses2-app/pkg/config"
)

func TestNewEmailMessage(t *testing.T) {
	tests := []struct {
		name         string
		emailConfig  config.EmailConfig
		to           []string
		templateData TemplateData
		expected     *EmailMessage
		hasError     bool
	}{
		{
			name: "Create email message",
			emailConfig: config.EmailConfig{
				From:    "test_from@example.com",
				Subject: "Test Subject",
				Body:    "The current exchange rate is {{.Rate}}.",
			},
			to: []string{"test_to@example.com"},
			templateData: TemplateData{
				Rate: "200",
			},
			expected: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "The current exchange rate is 200.",
			},
			hasError: false,
		},
		{
			name: "Bad template",
			emailConfig: config.EmailConfig{
				From:    "test_from@example.com",
				Subject: "Test Subject",
				Body:    "The current exchange rate is {{.Rate",
			},
			to: []string{"test_to@example.com"},
			templateData: TemplateData{
				Rate: "200",
			},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailMessage, err := NewEmailMessage(tt.emailConfig, tt.to, tt.templateData)
			if (err != nil) != tt.hasError {
				t.Errorf("NewEmailMessage() error = %v, wantErr %v", err, tt.hasError)
			}

			if err != nil {
				return
			}

			if emailMessage.from != tt.expected.from {
				t.Errorf("From: got %v, want %v", emailMessage.from, tt.expected.from)
			}

			if len(emailMessage.to) != len(tt.expected.to) || emailMessage.to[0] != tt.expected.to[0] {
				t.Errorf("To: got %v, want %v", emailMessage.to, tt.expected.to)
			}

			if emailMessage.subject != tt.expected.subject {
				t.Errorf("Subject: got %v, want %v", emailMessage.subject, tt.expected.subject)
			}

			if emailMessage.body != tt.expected.body {
				t.Errorf("Body: got %v, want %v", emailMessage.body, tt.expected.body)
			}
		})
	}
}

func TestPrepare(t *testing.T) {
	tests := []struct {
		name     string
		message  *EmailMessage
		expected string
	}{
		{
			name: "Prepare message",
			message: &EmailMessage{
				from:    "test_from@example.com",
				to:      []string{"test_to@example.com"},
				subject: "Test Subject",
				body:    "Test Body",
			},
			expected: "From: test_from@example.com\r\n" +
				"To: test_to@example.com\r\n" +
				"Subject: Test Subject\r\n" +
				"\r\nTest Body\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepared := tt.message.Prepare()

			if string(prepared) != tt.expected {
				t.Errorf("Prepared message: got %v, want %v", string(prepared), tt.expected)
			}
		})
	}
}
