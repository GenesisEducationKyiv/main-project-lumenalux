package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type StubExchangeRateService struct {
	rate float32
	err  error
}

func (m *StubExchangeRateService) ExchangeRate() (float32, error) {
	return m.rate, m.err
}

type StubEmailSubscriptionService struct {
	subscribeErr     error
	subscriptionsErr error
}

func (m *StubEmailSubscriptionService) Subscribe(email string) error {
	return m.subscribeErr
}

func (m *StubEmailSubscriptionService) Subscriptions() ([]string, error) {
	if m.subscriptionsErr != nil {
		return nil, m.subscriptionsErr
	}
	return []string{"subscriber1@example.com", "subscriber2@example.com"}, nil
}

type StubEmailSenderService struct {
	sendErr error
}

func (m *StubEmailSenderService) SendExchangeRate(
	rate float32,
	subscribers []string,
) error {
	return m.sendErr
}

func (m *StubEmailSubscriptionService) IsSubscribed(email string) (bool, error) {
	return true, nil
}

type StubExchangeRateServiceError struct{}

func (m *StubExchangeRateServiceError) ExchangeRate() (float32, error) {
	return 0, errors.New("Exchange rate error")
}

type StubEmailSubscriptionServiceError struct{}

func (m *StubEmailSubscriptionServiceError) Subscribe(email string) error {
	return errors.New("Subscription error")
}

func (m *StubEmailSubscriptionServiceError) IsSubscribed(email string) (bool, error) {
	return false, errors.New("Check subscription error")
}

func (m *StubEmailSubscriptionServiceError) Subscriptions() ([]string, error) {
	return nil, errors.New("Get subscriptions error")
}

type StubEmailSenderServiceError struct{}

func (m *StubEmailSenderServiceError) SendExchangeRate(
	rate float32,
	subscribers []string,
) error {
	return errors.New("Send email error")
}

func TestGetRate(t *testing.T) {
	tests := []struct {
		name           string
		exchangeRate   float32
		exchangeErr    error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Exchange rate",
			exchangeRate:   1.5,
			exchangeErr:    nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "1.5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				&StubExchangeRateService{rate: tt.exchangeRate, err: tt.exchangeErr},
				&StubEmailSubscriptionService{},
				&StubEmailSenderService{},
			)

			req, _ := http.NewRequest("GET", "/rate", nil)

			rr := httptest.NewRecorder()

			controller.GetRate(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("GetRate returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}

			actual := strings.TrimSpace(rr.Body.String())
			if actual != tt.expectedBody {
				t.Errorf("GetRate returned unexpected body: got %s, expected %s", actual, tt.expectedBody)
			}
		})
	}
}

func TestGetRateError(t *testing.T) {
	tests := []struct {
		name           string
		service        RateService
		expectedStatus int
	}{
		{
			name:           "Exchange rate error",
			service:        &StubExchangeRateServiceError{},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				tt.service,
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

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("GetRate returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSubscribeEmail(t *testing.T) {
	tests := []struct {
		name           string
		subscribeErr   error
		expectedStatus int
	}{
		{
			name:           "Happy path",
			subscribeErr:   nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				&StubExchangeRateService{},
				&StubEmailSubscriptionService{subscribeErr: tt.subscribeErr},
				&StubEmailSenderService{},
			)

			req, _ := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))

			rr := httptest.NewRecorder()

			controller.SubscribeEmail(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("SubscribeEmail returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSubscribeEmailError(t *testing.T) {
	tests := []struct {
		name           string
		service        SubscriptionService
		expectedStatus int
	}{
		{
			name:           "Subscription error",
			service:        &StubEmailSubscriptionServiceError{},
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				&StubExchangeRateService{},
				tt.service,
				&StubEmailSenderService{},
			)

			req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(controller.SubscribeEmail)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("SubscribeEmail returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSendEmails(t *testing.T) {
	tests := []struct {
		name             string
		exchangeRate     float32
		exchangeErr      error
		subscriptionsErr error
		sendErr          error
		expectedStatus   int
	}{
		{
			name:             "Send emails",
			exchangeRate:     1.5,
			exchangeErr:      nil,
			subscriptionsErr: nil,
			sendErr:          nil,
			expectedStatus:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				&StubExchangeRateService{rate: tt.exchangeRate, err: tt.exchangeErr},
				&StubEmailSubscriptionService{subscriptionsErr: tt.subscriptionsErr},
				&StubEmailSenderService{sendErr: tt.sendErr},
			)

			req, _ := http.NewRequest("POST", "/send", nil)

			rr := httptest.NewRecorder()

			controller.SendEmails(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("SendEmails returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestSendEmailsError(t *testing.T) {
	tests := []struct {
		name                string
		exchangeRateService RateService
		subscriptionService SubscriptionService
		emailSenderService  SenderService
		expectedStatus      int
	}{
		{
			name:                "Exchange rate error",
			exchangeRateService: &StubExchangeRateServiceError{},
			subscriptionService: &StubEmailSubscriptionService{},
			emailSenderService:  &StubEmailSenderService{},
			expectedStatus:      http.StatusBadRequest,
		},
		{
			name:                "Subscription service error",
			exchangeRateService: &StubExchangeRateService{},
			subscriptionService: &StubEmailSubscriptionServiceError{},
			emailSenderService:  &StubEmailSenderService{},
			expectedStatus:      http.StatusBadRequest,
		},
		{
			name:                "Email sender service error",
			exchangeRateService: &StubExchangeRateService{},
			subscriptionService: &StubEmailSubscriptionService{},
			emailSenderService:  &StubEmailSenderServiceError{},
			expectedStatus:      http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				tt.exchangeRateService,
				tt.subscriptionService,
				tt.emailSenderService,
			)

			req, err := http.NewRequest("POST", "/send", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(controller.SendEmails)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("SendEmails returned wrong status code: got %v, expected %v", status, tt.expectedStatus)
			}
		})
	}
}
