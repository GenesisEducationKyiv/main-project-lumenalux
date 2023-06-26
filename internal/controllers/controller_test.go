package controllers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubExchangeRateService struct{}

func (m *StubExchangeRateService) ExchangeRate() (float32, error) {
	return 1.5, nil
}

type StubEmailSubscriptionService struct{}

func (m *StubEmailSubscriptionService) Subscribe(email string) error {
	return nil
}

func (m *StubEmailSubscriptionService) IsSubscribed(email string) (bool, error) {
	return true, nil
}

func (m *StubEmailSubscriptionService) Subscriptions() ([]string, error) {
	return []string{"subscriber1@example.com", "subscriber2@example.com"}, nil
}

type StubEmailSenderService struct{}

func (m *StubEmailSenderService) SendExchangeRate(
	rate float32,
	subscribers []string,
) error {
	return nil
}

func TestGetRate(t *testing.T) {

	controller := NewAppController(
		&StubExchangeRateService{},
		&StubEmailSubscriptionService{},
		&StubEmailSenderService{},
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

func TestSubscribeEmail(t *testing.T) {

	controller := NewAppController(
		&StubExchangeRateService{},
		&StubEmailSubscriptionService{},
		&StubEmailSenderService{},
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

func TestSendEmails(t *testing.T) {

	controller := NewAppController(&StubExchangeRateService{}, &StubEmailSubscriptionService{}, &StubEmailSenderService{})

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
