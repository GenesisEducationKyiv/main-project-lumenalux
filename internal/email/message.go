package email

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"gses2-app/pkg/config"
)

const emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
{{.Body}}`

var (
	errParseTemplate  = errors.New("parse template error")
	errExecuteTempate = errors.New("cannot execute email")
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

func (e *EmailMessage) Prepare() ([]byte, error) {
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return nil, errParseTemplate
	}

	var message bytes.Buffer
	err = tmpl.Execute(&message, struct {
		From    string
		To      string
		Subject string
		Body    string
	}{
		From:    e.from,
		To:      strings.Join(e.to, ","),
		Subject: e.subject,
		Body:    e.body,
	})
	if err != nil {
		return nil, errExecuteTempate
	}

	return message.Bytes(), nil
}
