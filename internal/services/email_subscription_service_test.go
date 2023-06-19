package services

import (
	"os"
	"testing"
)

func TestSubscribe(t *testing.T) {
	testFilePath := "testEmails.csv"
	file, _ := os.Create(testFilePath)
	file.Close()

	defer os.Remove(testFilePath)

	service := NewEmailSubscriptionService(testFilePath)

	err := service.Subscribe("test@test.com")
	if err != nil {
		t.Errorf("Failed to subscribe: %v", err)
	}

	err = service.Subscribe("test@test.com")
	if err == nil {
		t.Errorf("Duplicate subscription didn't fail")
	}
}

func TestIsSubscribed(t *testing.T) {
	testFilePath := "testEmails.csv"
	file, _ := os.Create(testFilePath)
	file.Close()

	defer os.Remove(testFilePath)

	service := NewEmailSubscriptionService(testFilePath)

	subscribed, err := service.IsSubscribed("test@test.com")
	if subscribed || err != nil {
		t.Errorf("Non-subscribed email returned true")
	}

	service.Subscribe("test@test.com")
	subscribed, err = service.IsSubscribed("test@test.com")
	if !subscribed || err != nil {
		t.Errorf("Subscribed email returned false")
	}
}

func TestGetSubscriptions(t *testing.T) {
	testFilePath := "testEmails.csv"
	file, _ := os.Create(testFilePath)
	file.Close()

	defer os.Remove(testFilePath)

	service := NewEmailSubscriptionService(testFilePath)

	emails, err := service.Subscriptions()
	if err != nil || len(emails) != 0 {
		t.Errorf("GetSubscriptions didn't return empty list for empty file")
	}

	service.Subscribe("test@test.com")
	emails, err = service.Subscriptions()
	if err != nil || len(emails) != 1 || emails[0] != "test@test.com" {
		t.Errorf("GetSubscriptions didn't return correct emails")
	}
}
