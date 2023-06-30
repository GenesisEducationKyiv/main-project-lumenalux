// Subscription integration contains integration tests for the subscription layer,
// specifically testing the integration between the subscription service and
// the storage

package integration

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"gses2-app/internal/subscription"
	"gses2-app/pkg/storage"
)

type SubscribtionTest struct {
	Name           string
	Emails         []string
	Action         func(service *subscription.Service, emails []string) error
	ExpectedError  error
	ExpectedResult []string
}

func TestSubscriptionServiceIntegration(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	csvStorage := storage.NewCSVStorage(tmpFile.Name())
	service := subscription.NewService(csvStorage)

	tests := []SubscribtionTest{
		{
			Name:   "Subscribe a new email",
			Emails: []string{"test1@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				return service.Subscribe(emails[0])
			},
		},
		{
			Name:   "Check if an email is subscribed",
			Emails: []string{"test1@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				_, err := service.IsSubscribed(emails[0])
				return err
			},
			ExpectedResult: []string{"test1@example.com"},
		},
		{
			Name:   "Subscribe an already subscribed email",
			Emails: []string{"test1@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				return service.Subscribe(emails[0])
			},
			ExpectedError: subscription.ErrAlreadySubscribed,
		},
		{
			Name:   "Get all subscriptions",
			Emails: []string{},
			Action: func(service *subscription.Service, emails []string) error {
				_, err := service.Subscriptions()
				return err
			},
			ExpectedResult: []string{"test1@example.com"},
		},
		{
			Name:   "Subscribe multiple emails",
			Emails: []string{"test2@example.com", "test3@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				for _, email := range emails {
					if err := service.Subscribe(email); err != nil {
						return err
					}
				}
				return nil
			},
			ExpectedResult: []string{"test1@example.com", "test2@example.com", "test3@example.com"},
		},
		{
			Name:   "Subscribe new and already subscribed emails",
			Emails: []string{"test4@example.com", "test1@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				for _, email := range emails {
					err := service.Subscribe(email)
					if err != nil && !errors.Is(err, subscription.ErrAlreadySubscribed) {
						return err
					}
				}
				return nil
			},
			ExpectedResult: []string{"test1@example.com", "test2@example.com", "test3@example.com", "test4@example.com"},
		},
		{
			Name:   "Check if all emails are subscribed",
			Emails: []string{"test1@example.com", "test2@example.com", "test3@example.com", "test4@example.com"},
			Action: func(service *subscription.Service, emails []string) error {
				for _, email := range emails {
					subscribed, err := service.IsSubscribed(email)
					if err != nil || !subscribed {
						if err != nil {
							return err
						}
						return errors.New("Email not subscribed")
					}
				}
				return nil
			},
			ExpectedResult: []string{"test1@example.com", "test2@example.com", "test3@example.com", "test4@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			runTest(t, tt, service)
		})
	}
}

func runTest(t *testing.T, test SubscribtionTest, service *subscription.Service) {
	err := test.Action(service, test.Emails)
	checkError(t, err, test.ExpectedError)
	checkExpectedResult(t, service, test.ExpectedResult)
}

func checkError(t *testing.T, err error, expectedError error) {
	if !errors.Is(err, expectedError) {
		t.Fatalf("Expected error %v, got: %v", expectedError, err)
	}
}

func checkExpectedResult(t *testing.T, service *subscription.Service, expectedResult []string) {
	if expectedResult == nil {
		return
	}

	subscriptions, err := service.Subscriptions()
	if err != nil {
		t.Fatalf("Failed to get all subscriptions: %v", err)
	}

	if !cmp.Equal(subscriptions, expectedResult) {
		t.Errorf("Unexpected subscriptions. Got: %v, Expected: %v", subscriptions, expectedResult)
	}
}
