package email

import (
	"testing"

	"gses2-app/pkg/config"
)

func TestNewEmailMessage(t *testing.T) {
	emailConfig := config.EmailConfig{
		From:    "test_from@example.com",
		Subject: "Test Subject",
		Body:    "The current exchange rate is {{.Rate}}.",
	}
	to := []string{"test_to@example.com"}
	templateData := TemplateData{
		Rate: "200",
	}

	emailMessage, err := NewEmailMessage(emailConfig, to, templateData)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if emailMessage.from != emailConfig.From {
		t.Errorf("From: got %v, want %v", emailMessage.from, emailConfig.From)
	}

	if len(emailMessage.to) != len(to) || emailMessage.to[0] != to[0] {
		t.Errorf("To: got %v, want %v", emailMessage.to, to)
	}

	if emailMessage.subject != emailConfig.Subject {
		t.Errorf("Subject: got %v, want %v", emailMessage.subject, emailConfig.Subject)
	}

	expectedBody := "The current exchange rate is 200."
	if emailMessage.body != expectedBody {
		t.Errorf("Body: got %v, want %v", emailMessage.body, expectedBody)
	}
}

func TestPrepare(t *testing.T) {
	emailMessage := &EmailMessage{
		from:    "test_from@example.com",
		to:      []string{"test_to@example.com"},
		subject: "Test Subject",
		body:    "Test Body",
	}

	prepared := emailMessage.Prepare()

	expected := "From: test_from@example.com\r\n" +
		"To: test_to@example.com\r\n" +
		"Subject: Test Subject\r\n" +
		"\r\nTest Body\r\n"
	if string(prepared) != expected {
		t.Errorf("Prepared message: got %v, want %v", string(prepared), expected)
	}
}
