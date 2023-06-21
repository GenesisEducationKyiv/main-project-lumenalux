package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockExchangeRateService struct{}

func (m *MockExchangeRateService) ExchangeRate() (float32, error) {
	return 1.5, nil
}

type MockEmailSubscriptionService struct{}

func (m *MockEmailSubscriptionService) Subscribe(email string) error {
	return nil
}

func (m *MockEmailSubscriptionService) IsSubscribed(email string) (bool, error) {
	return true, nil
}

func (m *MockEmailSubscriptionService) Subscriptions() ([]string, error) {
	return []string{"subscriber1@example.com", "subscriber2@example.com"}, nil
}

type MockEmailSenderService struct{}

func (m *MockEmailSenderService) SendExchangeRate(
	rate float32,
	subscribers []string,
) error {
	return nil
}

func TestAppController_GetRate(t *testing.T) {

	controller := NewAppController(
		&MockExchangeRateService{},
		&MockEmailSubscriptionService{},
		&MockEmailSenderService{},
	)

	req, err := http.NewRequest("GET", "/rate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.GetRate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("GetRate returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}

	expected := "1.5"
	actual := strings.TrimSpace(rr.Body.String())
	if actual != expected {
		t.Errorf("GetRate returned unexpected body: got %s, expected %s", actual, expected)
	}
}

func TestAppController_SubscribeEmail(t *testing.T) {

	controller := NewAppController(
		&MockExchangeRateService{},
		&MockEmailSubscriptionService{},
		&MockEmailSenderService{},
	)

	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.SubscribeEmail)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("SubscribeEmail returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}

func TestAppController_SendEmails(t *testing.T) {

	controller := NewAppController(&MockExchangeRateService{}, &MockEmailSubscriptionService{}, &MockEmailSenderService{})

	req, err := http.NewRequest("POST", "/send", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(controller.SendEmails)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("SendEmails returned wrong status code: got %v, expected %v", status, http.StatusOK)
	}
}
