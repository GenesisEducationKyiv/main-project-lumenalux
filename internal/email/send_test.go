package email

import (
	"bytes"
	"errors"
	"testing"
)

func TestSendEmail(t *testing.T) {
	client := &StubSenderSMTPClient{}
	email := &EmailMessage{
		from:    "test_from@example.com",
		to:      []string{"test_to@example.com"},
		subject: "Test Subject",
		body:    "Test Body",
	}

	err := SendEmail(client, email)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if client.fromCalledWith != email.from {
		t.Errorf("Mail from: got %v, want %v", client.fromCalledWith, email.from)
	}

	if len(client.rcptCalledWith) != len(email.to) || client.rcptCalledWith[0] != email.to[0] {
		t.Errorf("Rcpt to: got %v, want %v", client.rcptCalledWith, email.to)
	}

	if !client.dataCalled {
		t.Error("Data was not called")
	}

	expectedMessage := email.Prepare()
	if !bytes.Equal(client.writeCalledWith, expectedMessage) {
		t.Errorf("Write called with: got %v, want %v", client.writeCalledWith, expectedMessage)
	}

	if !client.quitCalled {
		t.Error("Quit was not called")
	}
}

func TestSendEmailWriteError(t *testing.T) {
	client := &StubSenderSMTPClient{writeShouldReturn: errors.New("write error")}
	email := &EmailMessage{
		from:    "test_from@example.com",
		to:      []string{"test_to@example.com"},
		subject: "Test Subject",
		body:    "Test Body",
	}

	err := SendEmail(client, email)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	expectedError := "write error"
	if err.Error() != expectedError {
		t.Errorf("Error: got %v, want %v", err, expectedError)
	}
}
