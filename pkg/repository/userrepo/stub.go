package userrepo

import (
	"gses2-app/pkg/types"
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
