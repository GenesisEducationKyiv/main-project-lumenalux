package subscription

import (
	"gses2-app/pkg/repository/userrepo"
	"gses2-app/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestSubscription(t *testing.T) {
	t.Run("Subscribe", func(t *testing.T) {
		t.Parallel()

		subscriber := &types.User{Email: "test@example.com"}
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
			Users: []types.User{},
			Err:   userrepo.ErrAlreadyAdded,
		}
		service := NewService(userRepository)
		subscriber := &types.User{Email: "test@example.com"}

		err := service.Subscribe(subscriber)
		require.ErrorIs(
			t, err, ErrAlreadySubscribed,
			"expected error due to duplicate subscription",
		)
	})
}
