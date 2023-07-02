package sender

import (
	"fmt"

	"gses2-app/pkg/config"
)

type Service struct {
	config     *config.Config
	connection SMTPConnectionClient
}

func NewService(
	config *config.Config,
	dialer TLSConnectionDialer,
	factory SMTPClientFactory,
) (*Service, error) {
	client := NewSMTPClient(config.SMTP, dialer, factory)
	clientConnection, err := client.Connect()
	if err != nil {
		return nil, err
	}

	return &Service{config: config, connection: clientConnection}, nil
}

type TemplateData struct {
	Rate string
}

func (es *Service) SendExchangeRate(
	exchangeRate float32,
	emailAddresses []string,
) error {

	templateData := TemplateData{Rate: fmt.Sprintf("%.2f", exchangeRate)}
	emailMessage, err := NewEmailMessage(es.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	return SendEmail(es.connection, emailMessage)
}
