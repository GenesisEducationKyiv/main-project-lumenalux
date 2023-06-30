package subscription

import (
	"errors"
)

type Storage interface {
	Append(record ...string) error
	AllRecords() (records [][]string, err error)
}

type Service struct {
	Storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{Storage: storage}
}

var ErrAlreadySubscribed = errors.New("email is already subscribed")

func (s *Service) Subscribe(email string) error {
	subscribed, err := s.IsSubscribed(email)
	if err != nil {
		return err
	}
	if subscribed {
		return ErrAlreadySubscribed
	}

	return s.Storage.Append(email)
}

func (s *Service) IsSubscribed(email string) (bool, error) {
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

func (s *Service) Subscriptions() ([]string, error) {
	return s.allEmails()
}

func (s *Service) allEmails() ([]string, error) {
	records, err := s.Storage.AllRecords()
	if err != nil {
		return nil, err
	}

	return s.convertRecordsToEmails(records), nil
}

func (s *Service) convertRecordsToEmails(records [][]string) []string {
	emails := make([]string, len(records))
	for i, record := range records {
		emails[i] = record[0]
	}

	return emails
}
