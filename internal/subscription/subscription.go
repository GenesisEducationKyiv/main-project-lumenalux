package subscription

import (
	"errors"
	"gses2-app/pkg/repository/userrepo"
)

var (
	ErrAlreadySubscribed = errors.New("email is already subscribed")
	ErrUserRepository    = errors.New("user repository error")
)

type UserRepository interface {
	Add(user *userrepo.User) error
	All() ([]userrepo.User, error)
}

type Service struct {
	userRepository UserRepository
}

func NewService(userRepository UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Subscribe(user *userrepo.User) error {
	err := s.userRepository.Add(user)
	if errors.Is(err, userrepo.ErrAlreadyAdded) {
		return ErrAlreadySubscribed
	}

	if err != nil {
		return errors.Join(err, ErrUserRepository)
	}

	return nil
}

func (s *Service) Subscriptions() ([]userrepo.User, error) {
	return s.userRepository.All()
}
