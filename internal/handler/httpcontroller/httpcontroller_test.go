package httpcontroller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
	"gses2-app/internal/core/service/subscription"
)

var (
	errSubscriptions = errors.New("get subscriptions error")
	errExchangeRate  = errors.New("exchange rate error")
	errSendEmail     = errors.New("send email error")
)

type StubExchangeRateService struct {
	rate port.Rate
	err  error
}

func (m *StubExchangeRateService) ExchangeRate() (port.Rate, error) {
	return m.rate, m.err
}

type StubEmailSubscriptionService struct {
	subscribeErr     error
	subscriptions    []port.User
	subscriptionsErr error
	isSubscribedErr  error
}

func (m *StubEmailSubscriptionService) Subscribe(subscriber *port.User) error {
	return m.subscribeErr
}

func (m *StubEmailSubscriptionService) Subscriptions() ([]port.User, error) {
	if m.subscriptionsErr != nil {
		return nil, m.subscriptionsErr
	}
	return m.subscriptions, nil
}

func (m *StubEmailSubscriptionService) IsSubscribed(subscriber port.User) (bool, error) {
	return true, m.isSubscribedErr
}

type StubEmailSenderService struct {
	sendErr error
}

func (m *StubEmailSenderService) SendExchangeRate(
	rate port.Rate,
	subscribers ...port.User,
) error {
	return m.sendErr
}

func TestGetRate(t *testing.T) {
	tests := []struct {
		name           string
		service        *StubExchangeRateService
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Exchange rate",
			service:        &StubExchangeRateService{rate: 1.5},
			expectedStatus: http.StatusOK,
			expectedBody:   "1.5",
		},
		{
			name:           "Exchange rate error",
			service:        &StubExchangeRateService{err: errExchangeRate},
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

			req, err := http.NewRequest(http.MethodGet, "/rate", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(controller.GetRate)
			handler.ServeHTTP(rr, req)

			require.Equal(t,
				tt.expectedStatus,
				rr.Code,
				"GetRate returned wrong status code: got %v, expected %v",
				rr.Code,
				tt.expectedStatus,
			)

			if tt.expectedBody != "" {
				actual := strings.TrimSpace(rr.Body.String())
				require.Equal(
					t,
					tt.expectedBody,
					actual,
					"GetRate returned unexpected body: got %s, expected %s",
					rr.Code,
					tt.expectedStatus,
				)
			}
		})
	}
}

func TestSubscribeEmail(t *testing.T) {
	tests := []struct {
		name           string
		service        *StubEmailSubscriptionService
		expectedStatus int
	}{
		{
			name:           "Subscribe email",
			service:        &StubEmailSubscriptionService{},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Subscription error",
			service: &StubEmailSubscriptionService{
				subscribeErr: subscription.ErrAlreadySubscribed,
			},
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

			req, err := http.NewRequest(http.MethodPost, "/subscribe", strings.NewReader("email=test@example.com"))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(controller.SubscribeEmail)
			handler.ServeHTTP(rr, req)

			require.Equal(
				t,
				tt.expectedStatus,
				rr.Code,
				"SubscribeEmail returned wrong status code: got %v, expected %v",
				rr.Code,
				tt.expectedStatus,
			)
		})
	}
}

func TestSendEmails(t *testing.T) {
	tests := []struct {
		name                string
		exchangeRateService *StubExchangeRateService
		subscriptionService *StubEmailSubscriptionService
		emailSenderService  *StubEmailSenderService
		expectedStatus      int
	}{
		{
			name: "Send emails",
			exchangeRateService: &StubExchangeRateService{
				rate: 1.5,
			},
			subscriptionService: &StubEmailSubscriptionService{
				subscriptions: convertEmailsToUsers(
					[]string{"subscriber1@example.com", "subscriber2@example.com"},
				),
			},
			emailSenderService: &StubEmailSenderService{},
			expectedStatus:     http.StatusOK,
		},
		{
			name: "Exchange rate error",
			exchangeRateService: &StubExchangeRateService{
				err: errExchangeRate,
			},
			subscriptionService: &StubEmailSubscriptionService{},
			emailSenderService:  &StubEmailSenderService{},
			expectedStatus:      http.StatusBadRequest,
		},
		{
			name:                "Subscription service error",
			exchangeRateService: &StubExchangeRateService{},
			subscriptionService: &StubEmailSubscriptionService{
				subscriptionsErr: errSubscriptions,
			},
			emailSenderService: &StubEmailSenderService{},
			expectedStatus:     http.StatusInternalServerError,
		},
		{
			name:                "Email sender service error",
			exchangeRateService: &StubExchangeRateService{},
			subscriptionService: &StubEmailSubscriptionService{},
			emailSenderService: &StubEmailSenderService{
				sendErr: errSendEmail,
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller := NewAppController(
				tt.exchangeRateService,
				tt.subscriptionService,
				tt.emailSenderService,
			)

			req, err := http.NewRequest(http.MethodPost, "/send", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(controller.SendEmails)
			handler.ServeHTTP(rr, req)

			require.Equal(
				t,
				tt.expectedStatus,
				rr.Code,
				"SendEmails returned wrong status code: got %v, expected %v",
				rr.Code,
				tt.expectedStatus,
			)
		})
	}
}

func convertEmailsToUsers(emails []string) []port.User {
	users := make([]port.User, len(emails))

	for i, email := range emails {
		users[i] = port.User{Email: email}
	}

	return users
}
