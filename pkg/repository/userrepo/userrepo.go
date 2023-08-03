package userrepo

import (
	"errors"
	"gses2-app/pkg/types"
)

var (
	ErrAlreadyAdded      = errors.New("user is already added")
	ErrCannotFindByEmail = errors.New("cannot find user by email")
	ErrCannotLoadUsers   = errors.New("cannot load users")
)

type Storage interface {
	Append(record ...string) error
	AllRecords() (records [][]string, err error)
}

type UserRepository struct {
	storage Storage
}

func NewUserRepository(storage Storage) *UserRepository {
	return &UserRepository{
		storage: storage,
	}
}

func (ur *UserRepository) Add(user *types.User) error {
	_, err := ur.FindByEmail(user.Email)

	isUserFound := !errors.Is(err, ErrCannotFindByEmail)
	if isUserFound {
		return ErrAlreadyAdded
	}

	if err != nil && isUserFound {
		return err
	}

	return ur.storage.Append(user.Email)
}

func (ur *UserRepository) FindByEmail(email string) (*types.User, error) {
	records, err := ur.storage.AllRecords()
	if err != nil {
		return &types.User{}, err
	}

	for _, e := range records {
		if e[0] == email {
			return &types.User{Email: email}, nil
		}
	}

	return &types.User{}, ErrCannotFindByEmail
}

func (ur *UserRepository) All() ([]types.User, error) {
	records, err := ur.storage.AllRecords()
	if err != nil {
		return nil, errors.Join(err, ErrCannotLoadUsers)
	}

	users := make([]types.User, len(records))
	for i, record := range records {
		users[i].Email = record[0]
	}

	return users, nil
}
