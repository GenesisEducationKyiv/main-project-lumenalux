package send

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEmailMessage(t *testing.T) {
	tests := []struct {
		name         string
		emailConfig  EmailConfig
		to           []string
		templateData TemplateData
		expected     *EmailMessage
		hasError     bool
	}{
		{
			name: "Create email message",
			emailConfig: EmailConfig{
				From:    "test_from@example.com",
				Subject: "Test Subject",
				Body:    "The current exchange rate is {{.Rate}}.",
			},
			to: []string{"test_to@example.com"},
			templateData: TemplateData{
				Rate: "200",
			},
			expected: &EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "The current exchange rate is 200.",
			},
			hasError: false,
		},
		{
			name: "Bad template",
			emailConfig: EmailConfig{
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			emailMessage, err := NewEmailMessage(tt.emailConfig, tt.to, tt.templateData)
			require.Equal(
				t,
				tt.hasError,
				err != nil,
				"NewEmailMessage() error = %v, wantErr %v",
				err,
				tt.hasError,
			)

			if err != nil {
				return
			}

			require.Equal(
				t,
				tt.expected.From,
				emailMessage.From,
				"From: got %v, want %v",
				emailMessage.From,
				tt.expected.From,
			)

			require.Equal(
				t,
				tt.expected.To,
				emailMessage.To,
				"To: got %v, want %v",
				emailMessage.To,
				tt.expected.To,
			)

			require.Equal(
				t,
				tt.expected.Subject,
				emailMessage.Subject,
				"Subject: got %v, want %v",
				emailMessage.Subject,
				tt.expected.Subject,
			)

			require.Equal(
				t,
				tt.expected.Body,
				emailMessage.Body,
				"Body: got %v, want %v",
				emailMessage.Body,
				tt.expected.Body,
			)
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
			name: "Prepare single recipient message",
			message: &EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expected: `From: test_from@example.com
To: test_to@example.com
Subject: Test Subject

Test Body`,
		},
		{
			name: "Prepare multiple recipient message",
			message: &EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to1@example.com", "test_to2@example.com"},
				Subject: "Test Subject",
				Body:    "Test Body",
			},
			expected: `From: test_from@example.com
To: test_to1@example.com,test_to2@example.com
Subject: Test Subject

Test Body`,
		},
		{
			name: "Prepare message with no body",
			message: &EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "Test Subject",
				Body:    "",
			},
			expected: `From: test_from@example.com
To: test_to@example.com
Subject: Test Subject

`,
		},
		{
			name: "Prepare message with no subject",
			message: &EmailMessage{
				From:    "test_from@example.com",
				To:      []string{"test_to@example.com"},
				Subject: "",
				Body:    "Test Body",
			},
			expected: `From: test_from@example.com
To: test_to@example.com
Subject: 

Test Body`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			prepared, err := tt.message.Prepare()

			require.NoError(t, err, "Prepared message: want %v but got error %v", tt.expected, err)
			require.Equal(
				t,
				tt.expected,
				string(prepared),
				"Prepared message: got \n%v, want \n%v",
				string(prepared),
				tt.expected,
			)
		})
	}
}
