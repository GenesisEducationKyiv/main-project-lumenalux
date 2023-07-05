package email

import (
	"fmt"

	"gses2-app/internal/sender/provider/email/message"
	"gses2-app/internal/sender/provider/email/send"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/pkg/config"
	"gses2-app/pkg/types"
)

type Provider struct {
	config     *config.Config
	connection smtp.SMTPConnectionClient
}

func NewProvider(
	config *config.Config,
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
	rate types.Rate,
	subscribers []types.User,
) error {

	emailAddresses := convertSubscribersToEmails(subscribers)

	templateData := message.TemplateData{Rate: fmt.Sprintf("%.2f", rate)}
	emailMessage, err := message.NewEmailMessage(p.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	return send.SendEmail(p.connection, emailMessage)
}

func convertSubscribersToEmails(subscribers []types.User) []string {
	emails := make([]string, len(subscribers))

	for i, subscriber := range subscribers {
		emails[i] = string(subscriber)
	}

	return emails
}
