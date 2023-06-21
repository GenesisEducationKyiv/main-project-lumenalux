package email

import (
	"fmt"

	"gses2-app/pkg/config"
)

type SenderService interface {
	SendExchangeRate(float32, []string) error
}

type SenderServiceImpl struct {
	config        *config.Config
	dialer        TLSConnectionDialer
	clientFactory SMTPClientFactory
}

func NewSenderService(
	config *config.Config,
	dialer TLSConnectionDialer,
	factory SMTPClientFactory,
) *SenderServiceImpl {
	return &SenderServiceImpl{
		config:        config,
		dialer:        dialer,
		clientFactory: factory,
	}
}

type TemplateData struct {
	Rate string
}

func (es *SenderServiceImpl) SendExchangeRate(
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
