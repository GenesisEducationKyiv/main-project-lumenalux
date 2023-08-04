package subscription

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gses2-app/internal/core/port"
)

type StubUserRepository struct {
	Users []port.User
	Err   error
}

func (s *StubUserRepository) Add(user *port.User) error {
	s.Users = append(s.Users, *user)
	return s.Err
}

func (s *StubUserRepository) FindByEmail(email string) (*port.User, error) {
	return &s.Users[0], s.Err
}

func (s *StubUserRepository) All() ([]port.User, error) {
	return s.Users, s.Err
}

func TestSubscription(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		t.Parallel()

		subscriber := &port.User{Email: "test@example.com"}
		userRepository := &StubUserRepository{}
		service := NewService(userRepository)

		err := service.Subscribe(subscriber)
		require.NoError(t, err)

		subscribers, err := service.Subscriptions()
		require.NoError(t, err)

		require.Equal(
			t, 1, len(subscribers),
			"expected subscribers list to contain one subscriber",
		)

		require.Equal(
			t, *subscriber, subscribers[0],
			"expected subscribers list to contain the subscriber",
		)
	})

	t.Run("Already subscribed", func(t *testing.T) {
		t.Parallel()

		userRepository := &StubUserRepository{
			Users: []port.User{},
			Err:   port.ErrAlreadyAdded,
		}
		service := NewService(userRepository)
		subscriber := &port.User{Email: "test@example.com"}

		err := service.Subscribe(subscriber)
		require.ErrorIs(
			t, err, ErrAlreadySubscribed,
			"expected error due to duplicate subscription",
		)
	})
}
