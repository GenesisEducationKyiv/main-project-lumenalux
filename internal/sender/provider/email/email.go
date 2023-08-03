package email

import (
	"fmt"

	"gses2-app/internal/rate"
	"gses2-app/internal/sender/provider/email/message"
	"gses2-app/internal/sender/provider/email/send"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/internal/user/repository"
)

type EmailSenderConfig struct {
	SMTP  smtp.SMTPConfig
	Email message.EmailConfig
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
	rate rate.Rate,
	subscribers []repository.User,
) error {

	emailAddresses := convertUsersToEmails(subscribers)

	templateData := message.TemplateData{Rate: fmt.Sprintf("%.2f", rate)}
	emailMessage, err := message.NewEmailMessage(p.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	return send.SendEmail(p.connection, emailMessage)
}

func convertUsersToEmails(users []repository.User) []string {
	emails := make([]string, len(users))

	for i, user := range users {
		emails[i] = user.Email
	}

	return emails
}
