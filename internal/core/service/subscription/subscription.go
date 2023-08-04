package subscription

import (
	"errors"
	"gses2-app/internal/core/port"
)

var (
	ErrAlreadySubscribed = errors.New("email is already subscribed")
	ErrUserRepository    = errors.New("user repository error")
)

type UserRepository interface {
	Add(user *port.User) error
	All() ([]port.User, error)
}

type Service struct {
	userRepository UserRepository
}

func NewService(userRepository UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Subscribe(user *port.User) error {
	err := s.userRepository.Add(user)
	if errors.Is(err, port.ErrAlreadyAdded) {
		return ErrAlreadySubscribed
	}

	if err != nil {
		return errors.Join(err, ErrUserRepository)
	}

	return nil
}

func (s *Service) Subscriptions() ([]port.User, error) {
	return s.userRepository.All()
}
