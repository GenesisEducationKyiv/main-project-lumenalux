package services

import (
	"strings"
	"testing"

	"gses2-app/pkg/config"
)

func TestNewSMTPClient(t *testing.T) {
	config := config.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     465,
		User:     "test@example.com",
		Password: "password",
	}
	client := NewSMTPClient(config)

	if client.host != config.Host ||
		client.port != config.Port ||
		client.user != config.User ||
		client.password != config.Password {
		t.Errorf("NewSMTPClient() didn't correctly initialize SMTPClient")
	}
}

func TestNewEmailMessage(t *testing.T) {
	config := config.EmailConfig{
		From:    "test@example.com",
		Subject: "Test",
		Body:    "The current exchange rate is {{.Rate}}",
	}
	data := TemplateData{Rate: "123.45"}
	to := []string{"receiver@example.com"}

	email, err := NewEmailMessage(config, to, data)
	if err != nil {
		t.Errorf("Error creating new email message: %v", err)
	}

	if !strings.Contains(email.body, data.Rate) {
		t.Errorf("Expected email body to contain rate %s, got %s", data.Rate, email.body)
	}
}
