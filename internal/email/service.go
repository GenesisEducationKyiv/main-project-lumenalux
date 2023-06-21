package email

import (
	"fmt"
	"log"
	"net/http"

	"gses2-app/pkg/config"
)

type SenderService interface {
	SendExchangeRate(float32, []string) int
}

type SenderServiceImpl struct {
	dialer        TLSConnectionDialer
	clientFactory SMTPClientFactory
}

func NewSenderService(dialer TLSConnectionDialer, factory SMTPClientFactory) *SenderServiceImpl {
	return &SenderServiceImpl{
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
) int {
	config := config.Current()

	client := NewSMTPClient(config.SMTP, es.dialer, es.clientFactory)
	clientConnection, err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}

	templateData := TemplateData{Rate: fmt.Sprintf("%.2f", exchangeRate)}
	email, err := NewEmailMessage(config.Email, emailAddresses, templateData)
	if err != nil {
		log.Fatal(err)
	}

	err = SendEmail(clientConnection, email)
	if err != nil {
		log.Fatal(err)
	}

	return http.StatusOK
}
