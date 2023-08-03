package repository

import (
	"errors"
)

const _emailKey = "email"

var (
	ErrAlreadyAdded      = errors.New("user is already added")
	ErrCannotFindByEmail = errors.New("cannot find user by email")
	ErrCannotLoadUsers   = errors.New("cannot load users")
)

// Represents a User entity
type User struct {
	Email string
}

type Storage interface {
	Append(record map[string]string) error
	AllRecords() (records []map[string]string, err error)
}

type UserRepository struct {
	storage Storage
}

func NewUserRepository(storage Storage) *UserRepository {
	return &UserRepository{
		storage: storage,
	}
}

func (ur *UserRepository) Add(user *User) error {
	_, err := ur.FindByEmail(user.Email)

	isUserFound := !errors.Is(err, ErrCannotFindByEmail)
	if isUserFound {
		return ErrAlreadyAdded
	}

	if err != nil && isUserFound {
		return err
	}

	return ur.storage.Append(map[string]string{_emailKey: user.Email})
}

func (ur *UserRepository) FindByEmail(email string) (*User, error) {
	records, err := ur.storage.AllRecords()
	if err != nil {
		return &User{}, err
	}

	for _, e := range records {
		if e[_emailKey] == email {
			return &User{Email: email}, nil
		}
	}

	return &User{}, ErrCannotFindByEmail
}

func (ur *UserRepository) All() ([]User, error) {
	records, err := ur.storage.AllRecords()
	if err != nil {
		return nil, errors.Join(err, ErrCannotLoadUsers)
	}

	users := make([]User, len(records))
	for i, record := range records {
		users[i].Email = record[_emailKey]
	}

	return users, nil
}
