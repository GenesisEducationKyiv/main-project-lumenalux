package email

import (
	"fmt"

	"gses2-app/pkg/config"
)

type SenderService struct {
	config     *config.Config
	connection SMTPConnectionClient
}

func NewSenderService(
	config *config.Config,
	dialer TLSConnectionDialer,
	factory SMTPClientFactory,
) (*SenderService, error) {
	client := NewSMTPClient(config.SMTP, dialer, factory)
	clientConnection, err := client.Connect()
	if err != nil {
		return nil, err
	}

	SMTPClient := &SenderService{
		config:     config,
		connection: clientConnection,
	}

	return SMTPClient, nil
}

type TemplateData struct {
	Rate string
}

func (es *SenderService) SendExchangeRate(
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
