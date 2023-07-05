package subscription

import (
	"gses2-app/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

type StubStorage struct {
	Records [][]string
	Error   error
}

func (m *StubStorage) Append(record ...string) error {
	m.Records = append(m.Records, record)
	return m.Error
}

func (m *StubStorage) AllRecords() ([][]string, error) {
	return m.Records, m.Error
}

func TestSubscription(t *testing.T) {
	t.Run("Subscribe and check subscriptions", func(t *testing.T) {
		t.Parallel()

		stubStorage := &StubStorage{Records: [][]string{}}
		service := NewService(stubStorage)
		subscriber := types.User("test@example.com")

		err := service.Subscribe(subscriber)
		require.NoError(t, err)

		subscribed, err := service.IsSubscribed(subscriber)
		require.NoError(t, err)
		require.True(t, subscribed, "expected email to be subscribed")

		subscriptions, err := service.Subscriptions()
		require.NoError(t, err)

		require.Equal(
			t, 1, len(subscriptions),
			"expected subscription list to contain one email",
		)

		require.Equal(
			t, subscriber, subscriptions[0],
			"expected subscription list to contain the email",
		)
	})

	t.Run("Subscribe twice", func(t *testing.T) {
		t.Parallel()

		stubStorage := &StubStorage{Records: [][]string{}}
		service := NewService(stubStorage)
		subscriber := types.User("test@example.com")

		err := service.Subscribe(subscriber)
		require.NoError(t, err)

		err = service.Subscribe(subscriber)
		require.ErrorIs(
			t, err, ErrAlreadySubscribed,
			"expected error due to duplicate subscription",
		)
	})

	t.Run("Subscribed with non-existent email", func(t *testing.T) {
		t.Parallel()

		stubStorage := &StubStorage{Records: [][]string{}}
		service := NewService(stubStorage)
		subscriber := types.User("test@example.com")

		subscribed, err := service.IsSubscribed(subscriber)
		require.NoError(t, err)
		require.False(t, subscribed, "expected email not to be subscribed")
	})
}
