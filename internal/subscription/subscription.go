package subscription

import (
	"errors"
	"gses2-app/pkg/repository/userrepo"
	"gses2-app/pkg/types"
)

var (
	ErrAlreadySubscribed = errors.New("email is already subscribed")
	ErrUserRepository    = errors.New("user repository error")
)

type UserRepository interface {
	Add(user *types.User) error
	All() ([]types.User, error)
}

type Service struct {
	userRepository UserRepository
}

func NewService(userRepository UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Subscribe(user *types.User) error {
	err := s.userRepository.Add(user)
	if errors.Is(err, userrepo.ErrAlreadyAdded) {
		return ErrAlreadySubscribed
	}

	if err != nil {
		return errors.Join(err, ErrUserRepository)
	}

	return nil
}

func (s *Service) Subscriptions() ([]types.User, error) {
	return s.userRepository.All()
}
