package email

import (
	"fmt"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/sender/email/send"
	"gses2-app/internal/repository/sender/smtp"
)

type EmailSenderConfig struct {
	SMTP  smtp.SMTPConfig
	Email send.EmailConfig
}

type Provider struct {
	config     *EmailSenderConfig
	connection smtp.SMTPConnectionClient
}

func NewProvider(
	config *EmailSenderConfig,
	dialer smtp.TLSConnectionDialer,
	factory smtp.SMTPClientFactory,
) (*Provider, error) {
	client := smtp.NewSMTPClient(config.SMTP, dialer, factory)
	clientConnection, err := client.Connect()
	if err != nil {
		return nil, err
	}

	return &Provider{config: config, connection: clientConnection}, nil
}

func (p *Provider) SendExchangeRate(
	rate port.Rate,
	subscribers []port.User,
) error {

	emailAddresses := convertUsersToEmails(subscribers)

	templateData := send.TemplateData{Rate: fmt.Sprintf("%.2f", rate)}
	emailMessage, err := send.NewEmailMessage(p.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	return send.SendEmail(p.connection, emailMessage)
}

func convertUsersToEmails(users []port.User) []string {
	emails := make([]string, len(users))

	for i, user := range users {
		emails[i] = user.Email
	}

	return emails
}
