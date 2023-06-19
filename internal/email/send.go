package email

import (
	"net/smtp"
)

func setMail(client *smtp.Client, from string) error {
	return client.Mail(from)
}

func setRecipients(client *smtp.Client, to []string) error {
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}
	return nil
}

func writeAndClose(client *smtp.Client, message []byte) error {
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

func SendEmail(client *smtp.Client, email *EmailMessage) error {
	err := setMail(client, email.from)
	if err != nil {
		return err
	}

	err = setRecipients(client, email.to)
	if err != nil {
		return err
	}

	err = writeAndClose(client, email.Prepare())
	if err != nil {
		return err
	}

	err = client.Quit()
	return err
}
