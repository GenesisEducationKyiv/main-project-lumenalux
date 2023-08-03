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

	"gses2-app/internal/controller"
	"gses2-app/internal/rate"
	"gses2-app/internal/sender"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/internal/subscription"
	"gses2-app/internal/transport"
	"gses2-app/pkg/config"
	"gses2-app/pkg/repository/userrepo"
	"gses2-app/pkg/types"

	emailSenderProvider "gses2-app/internal/sender/provider/email"
)

const _configPrefix = "GSES2_APP"

type StubSenderProvider struct {
	Err error
}

func (tp *StubSenderProvider) SendExchangeRate(
	rate types.Rate,
	subscribers []types.User,
) error {
	return tp.Err
}

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

type StubRateProvider struct {
	Rate         types.Rate
	Error        error
	ProviderName string
}

func (m *StubRateProvider) ExchangeRate() (types.Rate, error) {
	return m.Rate, m.Error
}

func (m *StubRateProvider) Name() string {
	return m.ProviderName
}

type StubUserRepository struct {
	Users []types.User
	Err   error
}

func (s *StubUserRepository) Add(user *types.User) error {
	s.Users = append(s.Users, *user)
	return s.Err
}

func (s *StubUserRepository) FindByEmail(email string) (*types.User, error) {
	return &s.Users[0], s.Err
}

func (s *StubUserRepository) All() ([]types.User, error) {
	return s.Users, s.Err
}

var (
	errRateProviderAnavailable = errors.New("rate provider unavailable")
	errSendMessage             = errors.New("failed to send a message")
)

func TestAppControllerIntegration(t *testing.T) {
	config := initConfig(t)

	defaultEmailSenderService := sender.NewService(
		&StubSenderProvider{},
	)

	defaultRateService := rate.NewService(&StubRateProvider{Rate: 42})

	defaultSubscriptionService := subscription.NewService(
		&StubUserRepository{},
	)

	tests := []struct {
		name                string
		requestMethod       string
		requestURL          string
		requestBody         io.Reader
		expectedStatus      int
		senderService       *sender.Service
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
				&StubUserRepository{Err: userrepo.ErrAlreadyAdded},
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
				&StubRateProvider{
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
				&StubUserRepository{Err: userrepo.ErrCannotLoadUsers},
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
			senderService: initEmailSenderService(
				t,
				config,
				&smtp.StubDialer{},
				&smtp.StubSMTPClientFactory{
					Client: &smtp.StubSMTPClient{MailErr: errSendMessage},
				},
			),
			rateService: defaultRateService,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req, err := http.NewRequest(tt.requestMethod, tt.requestURL, tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			appController := controller.NewAppController(
				tt.rateService,
				tt.subscriptionService,
				tt.senderService,
			)

			if tt.requestMethod == http.MethodPost {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}

			rr := httptest.NewRecorder()

			router := transport.NewHTTPRouter(appController)
			mux := http.NewServeMux()
			router.RegisterRoutes(mux)

			mux.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func initConfig(t *testing.T) *config.Config {
	envVariables := map[string]string{
		"GSES2_APP_SMTP_HOST":             "test.server.com",
		"GSES2_APP_SMTP_USER":             "testuser",
		"GSES2_APP_SMTP_PORT":             "465",
		"GSES2_APP_SMTP_PASSWORD":         "testpassword",
		"GSES2_APP_EMAIL_FROM":            "no.reply@test.info.api",
		"GSES2_APP_EMAIL_SUBJECT":         "BTC to UAH exchange rate",
		"GSES2_APP_EMAIL_BODY":            "The BTC to UAH rate is {{.Rate}}",
		"GSES2_APP_STORAGE_PATH":          "./storage/storage.csv",
		"GSES2_APP_HTTP_PORT":             "8080",
		"GSES2_APP_HTTP_TIMEOUT":          "10s",
		"GSES2_APP_KUNA_API_URL":          "https://www.example.com",
		"GSES2_APP_KUNA_API_DEFAULT_RATE": "0",
	}

	for key, value := range envVariables {
		t.Setenv(key, value)
	}

	config, err := config.Load(_configPrefix)
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	return &config
}

func initEmailSenderService(
	t *testing.T,
	config *config.Config,
	dialer smtp.TLSConnectionDialer,
	factory smtp.SMTPClientFactory,
) *sender.Service {
	provider, err := emailSenderProvider.NewProvider(config, dialer, factory)

	if err != nil {
		t.Fatalf("error creating email sender provider: %v", err)
	}

	return sender.NewService(provider)
}
