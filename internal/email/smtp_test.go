package email

import (
	"testing"

	"gses2-app/pkg/config"
)

func TestConnect(t *testing.T) {
	config := config.SMTPConfig{
		Host:     "smtp.example.com",
		Port:     587,
		User:     "user@example.com",
		Password: "password",
	}

	factory := MockSMTPClientFactory{Client: &MockSMTPClient{}}
	client := NewSMTPClient(config, &MockDialer{}, factory)

	smtpClient, err := client.Connect()

	if smtpClient == nil {
		t.Errorf("Expected smtp.Client, got nil")
	}

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
