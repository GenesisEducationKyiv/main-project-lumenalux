package subscription

import (
	"testing"

	"gses2-app/pkg/user/repository"

	"github.com/stretchr/testify/require"
)

type StubUserRepository struct {
	Users []repository.User
	Err   error
}

func (s *StubUserRepository) Add(user *repository.User) error {
	s.Users = append(s.Users, *user)
	return s.Err
}

func (s *StubUserRepository) FindByEmail(email string) (*repository.User, error) {
	return &s.Users[0], s.Err
}

func (s *StubUserRepository) All() ([]repository.User, error) {
	return s.Users, s.Err
}

func TestSubscription(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		t.Parallel()

		subscriber := &repository.User{Email: "test@example.com"}
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
			Users: []repository.User{},
			Err:   repository.ErrAlreadyAdded,
		}
		service := NewService(userRepository)
		subscriber := &repository.User{Email: "test@example.com"}

		err := service.Subscribe(subscriber)
		require.ErrorIs(
			t, err, ErrAlreadySubscribed,
			"expected error due to duplicate subscription",
		)
	})
}
