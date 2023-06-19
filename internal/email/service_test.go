package email

import (
	"net/http"
	"testing"
)

func TestSendExchangeRate(t *testing.T) {
	dialer := &MockDialer{}
	factory := &MockSMTPClientFactory{Client: &MockSMTPClient{}}

	service := NewSenderService(dialer, factory)

	t.Run("Successful SendExchangeRate", func(t *testing.T) {
		emailAddresses := []string{"test@example.com"}
		exchangeRate := float32(10.5)
		statusCode := service.SendExchangeRate(exchangeRate, emailAddresses)

		expectedStatusCode := http.StatusOK
		if statusCode != expectedStatusCode {
			t.Errorf("Expected status code %d, but got %d", expectedStatusCode, statusCode)
		}
	})
}
