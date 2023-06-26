package integration

import (
	"bytes"
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

func (s *StubStorage) Append(record []string) error {
	return s.err
}

func (s *StubStorage) AllRecords() (records [][]string, err error) {
	return s.records, s.err
}

func TestAppController_Integration(t *testing.T) {
	config, err := config.Load("../../config.yaml")
	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	emailSenderService, err := email.NewSenderService(
		&config,
		&email.StubDialer{Err: nil},
		&email.StubSMTPClientFactory{Client: &email.StubSMTPClient{}, Err: nil},
	)

	if err != nil {
		t.Fatalf("error creating email sender service: %v", err)
	}

	rateService := rate.NewService(rate.NewKunaProvider(config.KunaAPI, &http.Client{}))

	tests := []struct {
		name           string
		requestMethod  string
		requestURL     string
		requestBody    io.Reader
		expectedStatus int
		storageRecords [][]string
	}{
		{
			name:           "GetRate_OK",
			requestMethod:  "GET",
			requestURL:     "/api/rate",
			requestBody:    nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SubscribeEmail_OK",
			requestMethod:  "POST",
			requestURL:     "/api/subscribe",
			requestBody:    bytes.NewBufferString("email=test@test.com"),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SubscribeEmail_StatusConflict",
			requestMethod:  "POST",
			requestURL:     "/api/subscribe",
			requestBody:    bytes.NewBufferString("email=test@test.com"),
			expectedStatus: http.StatusConflict,
			storageRecords: [][]string{{"test@test.com"}},
		},
		{
			name:           "SendEmails_OK",
			requestMethod:  "POST",
			requestURL:     "/api/sendEmails",
			requestBody:    nil,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.requestMethod, tt.requestURL, tt.requestBody)
			if err != nil {
				t.Fatal(err)
			}

			emailSubscriptionService := subscription.NewService(
				&StubStorage{records: tt.storageRecords},
			)

			controller := controllers.NewAppController(
				rateService,
				emailSubscriptionService,
				emailSenderService,
			)

			if tt.requestMethod == "POST" {
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
