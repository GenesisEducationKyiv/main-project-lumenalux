package subscription

import (
	"errors"
	"gses2-app/pkg/types"
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

func (s *Service) Subscribe(subscriber types.Subscriber) error {
	subscribed, err := s.IsSubscribed(subscriber)
	if err != nil {
		return err
	}
	if subscribed {
		return ErrAlreadySubscribed
	}

	return s.Storage.Append(string(subscriber))
}

func (s *Service) IsSubscribed(subscriber types.Subscriber) (bool, error) {
	subscribers, err := s.allSubscribers()
	if err != nil {
		return false, err
	}

	for _, s := range subscribers {
		if s == subscriber {
			return true, nil
		}
	}

	return false, nil
}

func (s *Service) Subscriptions() ([]types.Subscriber, error) {
	return s.allSubscribers()
}

func (s *Service) allSubscribers() ([]types.Subscriber, error) {
	records, err := s.Storage.AllRecords()
	if err != nil {
		return nil, err
	}

	return s.convertRecordsToSubscribers(records), nil
}

func (s *Service) convertRecordsToSubscribers(records [][]string) []types.Subscriber {
	subscribers := make([]types.Subscriber, len(records))
	for i, record := range records {
		subscribers[i] = types.Subscriber(record[0])
	}

	return subscribers
}
