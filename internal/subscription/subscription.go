package subscription

import (
	"errors"
	"gses2-app/pkg/user/repository"
)

var (
	ErrAlreadySubscribed = errors.New("email is already subscribed")
	ErrUserRepository    = errors.New("user repository error")
)

type UserRepository interface {
	Add(user *repository.User) error
	All() ([]repository.User, error)
}

type Service struct {
	userRepository UserRepository
}

func NewService(userRepository UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) Subscribe(user *repository.User) error {
	err := s.userRepository.Add(user)
	if errors.Is(err, repository.ErrAlreadyAdded) {
		return ErrAlreadySubscribed
	}

	if err != nil {
		return errors.Join(err, ErrUserRepository)
	}

	return nil
}

func (s *Service) Subscriptions() ([]repository.User, error) {
	return s.userRepository.All()
}
