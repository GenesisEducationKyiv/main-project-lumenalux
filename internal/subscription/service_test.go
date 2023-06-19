package subscription

import (
	"errors"
	"testing"
)

type MockStorage struct {
	Records [][]string
	Error   error
}

func (m *MockStorage) Append(record []string) error {
	m.Records = append(m.Records, record)
	return m.Error
}

func (m *MockStorage) Read() ([][]string, error) {
	return m.Records, m.Error
}

func TestSubscription(t *testing.T) {
	t.Run("Subscribe and check subscriptions", func(t *testing.T) {
		mockStorage := &MockStorage{Records: [][]string{}}
		service := NewService(mockStorage)
		email := "test@example.com"

		err := service.Subscribe(email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		subscribed, err := service.IsSubscribed(email)
		if err != nil || !subscribed {
			t.Errorf("expected email to be subscribed, got: %v", err)
		}

		subscriptions, err := service.Subscriptions()
		if err != nil || len(subscriptions) != 1 || subscriptions[0] != email {
			t.Errorf("expected subscription list to contain email, got: %v", subscriptions)
		}
	})

	t.Run("Subscribe twice", func(t *testing.T) {
		mockStorage := &MockStorage{Records: [][]string{}}
		service := NewService(mockStorage)
		email := "test@example.com"

		err := service.Subscribe(email)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		err = service.Subscribe(email)
		if !errors.Is(err, ErrAlreadySubscribed) {
			t.Errorf("expected error due to duplicate subscription, got: %v", err)
		}
	})
}
