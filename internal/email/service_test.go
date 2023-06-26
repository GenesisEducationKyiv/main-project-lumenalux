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

	service, err := NewSenderService(config, dialer, factory)

	if err != nil {
		t.Errorf("Expected service, but got error %s", err)
	}

	t.Run("Successful SendExchangeRate", func(t *testing.T) {
		emailAddresses := []string{"test@example.com"}
		expectedStatusCode := http.StatusOK
		exchangeRate := float32(10.5)

		err := service.SendExchangeRate(exchangeRate, emailAddresses)
		if err != nil {
			t.Errorf("Expected status code %d, but got error %s", expectedStatusCode, err)
		}
	})
}
