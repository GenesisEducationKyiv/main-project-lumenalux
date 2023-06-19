package email

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"gses2-app/pkg/config"
)

type EmailMessage struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewEmailMessage(
	config config.EmailConfig,
	to []string,
	data TemplateData,
) (*EmailMessage, error) {
	tmpl, err := template.New("email").Parse(config.Body)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return nil, err
	}

	return &EmailMessage{
		from:    config.From,
		to:      to,
		subject: config.Subject,
		body:    body.String(),
	}, nil
}

func (e *EmailMessage) Prepare() []byte {
	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n%s\r\n",
		e.from, strings.Join(e.to, ","), e.subject, e.body)

	return []byte(message)
}
