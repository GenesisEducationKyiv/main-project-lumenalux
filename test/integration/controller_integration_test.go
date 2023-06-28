// Controller integration contains integration tests for the controller,
// covering the interactions between the email, rate, subscription, and
// transport layers.

package integration

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gses2-app/internal/controllers"
	"gses2-app/internal/email"
	"gses2-app/internal/rate"
	"gses2-app/internal/subscription"
	"gses2-app/internal/transport"
	"gses2-app/pkg/config"
)

type StubStorage struct {
	err     error
	records [][]string
}

func (s *StubStorage) Append(record ...string) error {
	return s.err
}

func (s *StubStorage) AllRecords() (records [][]string, err error) {
	return s.records, s.err
}

var (
	errRateProviderAnavailable = errors.New("rate provider unavailable")
	errSendMessage             = errors.New("failed to send a message")
	errGetSubscribtions        = errors.New("cannot get subscribtions")
)

func TestAppController_Integration(t *testing.T) {
	config, err := config.Load("../../configs/config.yaml")
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	defaultEmailSenderService := initSenderService(
		t,
		&config,
		&email.StubDialer{},
		&email.StubSMTPClientFactory{Client: &email.StubSMTPClient{}},
	)

	defaultRateService := rate.NewService(&rate.StubProvider{Rate: 42})
	defaultSubscriptionService := subscription.NewService(&StubStorage{})

	tests := []struct {
		name                string
		requestMethod       string
		requestURL          string
		requestBody         io.Reader
		expectedStatus      int
		senderService       *email.SenderService
		subscriptionService *subscription.Service
		rateService         *rate.Service
	}{
		{
			name:                "GetRate OK",
			requestMethod:       http.MethodGet,
			requestURL:          "/api/rate",
			requestBody:         nil,
			expectedStatus:      http.StatusOK,
			subscriptionService: defaultSubscriptionService,
			senderService:       defaultEmailSenderService,
			rateService:         defaultRateService,
		},
		{
			name:                "SubscribeEmail OK",
			requestMethod:       http.MethodPost,
			requestURL:          "/api/subscribe",
			requestBody:         bytes.NewBufferString("email=test@test.com"),
			expectedStatus:      http.StatusOK,
			senderService:       defaultEmailSenderService,
			subscriptionService: defaultSubscriptionService,
			rateService:         defaultRateService,
		},
		{
			name:           "SubscribeEmail StatusConflict",
			requestMethod:  http.MethodPost,
			requestURL:     "/api/subscribe",
			requestBody:    bytes.NewBufferString("email=test@test.com"),
			expectedStatus: http.StatusConflict,
			subscriptionService: subscription.NewService(
				&StubStorage{records: [][]string{{"test@test.com"}}},
			),
			senderService: defaultEmailSenderService,
			rateService:   defaultRateService,
		},
		{
			name:                "SendEmails OK",
			requestMethod:       http.MethodPost,
			requestURL:          "/api/sendEmails",
			requestBody:         nil,
			expectedStatus:      http.StatusOK,
			subscriptionService: defaultSubscriptionService,
			senderService:       defaultEmailSenderService,
			rateService:         defaultRateService,
		},
		{
			name:                "SendEmails BadRequest Rate Provider Unavailable",
			requestMethod:       http.MethodPost,
			requestURL:          "/api/sendEmails",
			requestBody:         nil,
			expectedStatus:      http.StatusBadRequest,
			subscriptionService: defaultSubscriptionService,
			senderService:       defaultEmailSenderService,
			rateService: rate.NewService(
				&rate.StubProvider{
					Error: errRateProviderAnavailable,
				},
			),
		},
		{
			name:           "SendEmails InternalServerError Subscribtions Error",
			requestMethod:  http.MethodPost,
			requestURL:     "/api/sendEmails",
			requestBody:    nil,
			expectedStatus: http.StatusInternalServerError,
			subscriptionService: subscription.NewService(
				&StubStorage{err: errGetSubscribtions},
			),
			senderService: defaultEmailSenderService,
			rateService:   defaultRateService,
		},
		{
			name:                "SendEmails InternalServerError Send Error",
			requestMethod:       http.MethodPost,
			requestURL:          "/api/sendEmails",
			requestBody:         nil,
			expectedStatus:      http.StatusInternalServerError,
			subscriptionService: defaultSubscriptionService,
			senderService: initSenderService(
				t,
				&config,
				&email.StubDialer{},
				&email.StubSMTPClientFactory{
					Client: &email.StubSMTPClient{MailErr: errSendMessage},
				},
			),
			rateService: defaultRateService,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.requestMethod, tt.requestURL, tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			controller := controllers.NewAppController(
				tt.rateService,
				tt.subscriptionService,
				tt.senderService,
			)

			if tt.requestMethod == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			rr := httptest.NewRecorder()

			router := transport.NewHTTPRouter(controller)
			mux := http.NewServeMux()
			router.RegisterRoutes(mux)

			mux.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func initSenderService(
	t *testing.T,
	config *config.Config,
	dialer email.TLSConnectionDialer,
	factory email.SMTPClientFactory,
) *email.SenderService {
	service, err := email.NewSenderService(config, dialer, factory)

	if err != nil {
		t.Fatalf("error creating email sender service: %v", err)
	}

	return service
}
