// Subscription integration contains integration tests for the subscription layer,
// specifically testing the integration between the subscription service and
// the storage

package integration

import (
	"errors"
	"os"
	"testing"

	"gses2-app/internal/subscription"
	"gses2-app/pkg/storage"
)

func TestSubscriptionServiceIntegration(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}

	defer os.Remove(tmpFile.Name())

	csvStorage := storage.NewCSVStorage(tmpFile.Name())

	service := subscription.NewService(csvStorage)

	t.Run("Subscribe a new email", func(t *testing.T) {
		email := "test@example.com"
		err := service.Subscribe(email)
		if err != nil {
			t.Fatalf("Failed to subscribe a new email: %v", err)
		}
	})

	t.Run("Check if an email is subscribed", func(t *testing.T) {
		email := "test@example.com"
		subscribed, err := service.IsSubscribed(email)
		if err != nil {
			t.Fatalf("Failed to check if an email is subscribed: %v", err)
		}
		if !subscribed {
			t.Fatal("Expected the email to be subscribed")
		}
	})

	t.Run("Subscribe an already subscribed email", func(t *testing.T) {
		email := "test@example.com"
		err := service.Subscribe(email)
		if !errors.Is(err, subscription.ErrAlreadySubscribed) {
			t.Fatalf("Expected ErrAlreadySubscribed, got: %v", err)
		}
	})

	t.Run("Get all subscriptions", func(t *testing.T) {
		email := "test@example.com"
		subscriptions, err := service.Subscriptions()
		if err != nil {
			t.Fatalf("Failed to get all subscriptions: %v", err)
		}
		if len(subscriptions) != 1 || subscriptions[0] != email {
			t.Fatalf("Unexpected subscriptions. Got: %v, Expected: [%v]", subscriptions, email)
		}
	})
}
