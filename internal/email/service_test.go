package email

import (
	"gses2-app/pkg/config"
	"net/http"
	"testing"
)

func TestSendExchangeRate(t *testing.T) {
	config := &config.Config{}
	dialer := &MockDialer{}
	factory := &MockSMTPClientFactory{Client: &MockSMTPClient{}}

	service := NewSenderService(config, dialer, factory)

	t.Run("Successful SendExchangeRate", func(t *testing.T) {
		emailAddresses := []string{"test@example.com"}
		expectedStatusCode := http.StatusOK
		exchangeRate := float32(10.5)

		statusCode, err := service.SendExchangeRate(exchangeRate, emailAddresses)
		if err != nil {
			t.Errorf("Expected status code %d, but got error %s", expectedStatusCode, err)
		}

		if statusCode != expectedStatusCode {
			t.Errorf("Expected status code %d, but got %d", expectedStatusCode, statusCode)
		}
	})
}
