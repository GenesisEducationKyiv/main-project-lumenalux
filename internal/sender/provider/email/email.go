package email

import (
	"fmt"

	"gses2-app/internal/sender/provider/email/message"
	"gses2-app/internal/sender/provider/email/send"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/pkg/config"
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
	exchangeRate float32,
	emailAddresses []string,
) error {

	templateData := message.TemplateData{Rate: fmt.Sprintf("%.2f", exchangeRate)}
	emailMessage, err := message.NewEmailMessage(p.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	return send.SendEmail(p.connection, emailMessage)
}
