package subscription

import (
	"errors"
)

type Storage interface {
	Append(record []string) error
	AllRecords() (records [][]string, err error)
}

type Service interface {
	Subscribe(email string) error
	IsSubscribed(email string) (bool, error)
	Subscriptions() ([]string, error)
}

type ServiceImpl struct {
	Storage Storage
}

func NewService(storage Storage) *ServiceImpl {
	return &ServiceImpl{Storage: storage}
}

var ErrAlreadySubscribed = errors.New("email is already subscribed")

func (s *ServiceImpl) Subscribe(email string) error {
	subscribed, err := s.IsSubscribed(email)
	if err != nil {
		return err
	}
	if subscribed {
		return ErrAlreadySubscribed
	}

	return s.Storage.Append([]string{email})
}

func (s *ServiceImpl) IsSubscribed(email string) (bool, error) {
	emails, err := s.allEmails()
	if err != nil {
		return false, err
	}

	for _, e := range emails {
		if e == email {
			return true, nil
		}
	}

	return false, nil
}

func (s *ServiceImpl) Subscriptions() ([]string, error) {
	return s.allEmails()
}

func (s *ServiceImpl) allEmails() ([]string, error) {
	records, err := s.Storage.AllRecords()
	if err != nil {
		return nil, err
	}

	return s.convertRecordsToEmails(records), nil
}

func (s *ServiceImpl) convertRecordsToEmails(records [][]string) []string {
	emails := make([]string, len(records))
	for i, record := range records {
		emails[i] = record[0]
	}

	return emails
}
