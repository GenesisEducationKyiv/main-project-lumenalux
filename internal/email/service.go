package email

import (
	"fmt"

	"gses2-app/pkg/config"
)

type SenderService struct {
	config        *config.Config
	dialer        TLSConnectionDialer
	clientFactory SMTPClientFactory
}

func NewSenderService(
	config *config.Config,
	dialer TLSConnectionDialer,
	factory SMTPClientFactory,
) *SenderService {
	return &SenderService{
		config:        config,
		dialer:        dialer,
		clientFactory: factory,
	}
}

type TemplateData struct {
	Rate string
}

func (es *SenderService) SendExchangeRate(
	exchangeRate float32,
	emailAddresses []string,
) error {

	client := NewSMTPClient(es.config.SMTP, es.dialer, es.clientFactory)
	clientConnection, err := client.Connect()
	if err != nil {
		return err
	}

	templateData := TemplateData{Rate: fmt.Sprintf("%.2f", exchangeRate)}
	email, err := NewEmailMessage(es.config.Email, emailAddresses, templateData)
	if err != nil {
		return err
	}

	err = SendEmail(clientConnection, email)
	if err != nil {
		return err
	}

	return nil
}
