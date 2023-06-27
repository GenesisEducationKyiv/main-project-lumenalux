package email

import (
	"errors"
	"io"
)

var (
	errNoRecepients = errors.New("no recepiets")
)

type MailClient interface {
	Mail(string) error
	Rcpt(string) error
	Data() (io.WriteCloser, error)
	Quit() error
}

func setMail(client MailClient, from string) error {
	return client.Mail(from)
}

func setRecipients(client MailClient, to []string) error {
	if len(to) == 0 {
		return errNoRecepients
	}

	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}

	return nil
}

func writeAndClose(client MailClient, message []byte) error {
	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	if err != nil {
		return err
	}

	err = writer.Close()
	return err
}

func SendEmail(client MailClient, email *EmailMessage) error {
	err := setMail(client, email.from)
	if err != nil {
		return err
	}

	err = setRecipients(client, email.to)
	if errors.Is(err, errNoRecepients) {
		return nil
	}

	if err != nil {
		return err
	}

	return writeAndClose(client, email.Prepare())
}
