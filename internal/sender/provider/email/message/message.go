package message

import (
	"bytes"
	"errors"
	"strings"
	"text/template"

	"gses2-app/pkg/config"
)

const _emailTemplate = `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}

{{.Body}}`

var (
	errParseTemplate   = errors.New("parse template error")
	errExecuteTemplate = errors.New("cannot execute email")
)

type TemplateData struct {
	Rate string
}

type EmailMessage struct {
	From    string
	To      []string
	Subject string
	Body    string
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
		From:    config.From,
		To:      to,
		Subject: config.Subject,
		Body:    body.String(),
	}, nil
}

func (e *EmailMessage) Prepare() ([]byte, error) {
	tmpl, err := template.New("email").Parse(_emailTemplate)
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
		From:    e.From,
		To:      strings.Join(e.To, ","),
		Subject: e.Subject,
		Body:    e.Body,
	})
	if err != nil {
		return nil, errExecuteTemplate
	}

	return message.Bytes(), nil
}
