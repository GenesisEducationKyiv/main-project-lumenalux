package services

import (
	"encoding/csv"
	"errors"
	"os"
)

type EmailSubscriptionService interface {
	Subscribe(email string) error
	IsSubscribed(email string) (bool, error)
	Subscriptions() ([]string, error)
}

type EmailSubscriptionServiceImpl struct {
	FilePath string
}

func NewEmailSubscriptionService(filePath string) EmailSubscriptionServiceImpl {
	return EmailSubscriptionServiceImpl{FilePath: filePath}
}

func (s *EmailSubscriptionServiceImpl) Subscribe(email string) error {

	subscribed, err := s.IsSubscribed(email)
	if err != nil {
		return err
	}
	if subscribed {
		return errors.New("email is already subscribed")
	}

	return s.appendEmailToFile(email)
}

func (s *EmailSubscriptionServiceImpl) IsSubscribed(email string) (bool, error) {
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

func (s *EmailSubscriptionServiceImpl) Subscriptions() ([]string, error) {
	return s.allEmails()
}

func (s *EmailSubscriptionServiceImpl) allEmails() ([]string, error) {
	records, err := s.readCSVFile()
	if err != nil {
		return nil, err
	}

	return s.convertRecordsToEmails(records), nil
}

func (s *EmailSubscriptionServiceImpl) readCSVFile() ([][]string, error) {

	f, err := os.Open(s.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	return r.ReadAll()
}

func (s *EmailSubscriptionServiceImpl) convertRecordsToEmails(records [][]string) []string {
	emails := make([]string, len(records))
	for i, record := range records {
		emails[i] = record[0]
	}

	return emails
}

func (s *EmailSubscriptionServiceImpl) appendEmailToFile(email string) error {

	f, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err = w.Write([]string{email}); err != nil {
		return err
	}
	w.Flush()

	return w.Error()
}
