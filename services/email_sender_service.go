package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

type EmailSenderService interface {
	SendExchangeRate(float32, []string) int
}

type EmailSenderServiceImpl struct {
	ConfigFilePath string
}

func NewEmailSenderService(configFilePath string) *EmailSenderServiceImpl {
	return &EmailSenderServiceImpl{ConfigFilePath: configFilePath}
}

type TemplateData struct {
	Rate string
}

type Config struct {
	SMTP  SMTPConfig  `yaml:"smtp"`
	Email EmailConfig `yaml:"email"`
}

type EmailConfig struct {
	From    string `yaml:"from"`
	Subject string `yaml:"subject"`
	Body    string `yaml:"body"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

func (c *Config) loadFromYamlFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, c)
	return err
}

type SMTPClient struct {
	host     string
	port     int
	user     string
	password string
}

func NewSMTPClient(config SMTPConfig) *SMTPClient {
	return &SMTPClient{
		host:     config.Host,
		port:     config.Port,
		user:     config.User,
		password: config.Password,
	}
}

func (c *SMTPClient) createTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.host,
	}
}

func (c *SMTPClient) createConnection(tlsConfig *tls.Config) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", c.host+":"+strconv.Itoa(c.port), tlsConfig)
	return conn, err
}

func (c *SMTPClient) createSMTPClient(conn *tls.Conn) (*smtp.Client, error) {
	client, err := smtp.NewClient(conn, c.host)
	return client, err
}

func (c *SMTPClient) authenticate(client *smtp.Client) error {
	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	return client.Auth(auth)
}

func (c *SMTPClient) Connect() (*smtp.Client, error) {
	tlsConfig := c.createTLSConfig()
	conn, err := c.createConnection(tlsConfig)
	if err != nil {
		return nil, err
	}

	client, err := c.createSMTPClient(conn)
	if err != nil {
		return nil, err
	}

	err = c.authenticate(client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

type EmailMessage struct {
	from    string
	to      []string
	subject string
	body    string
}

func NewEmailMessage(config EmailConfig, to []string, data TemplateData) (*EmailMessage, error) {
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

func setMail(client *smtp.Client, from string) error {
	return client.Mail(from)
}

func setRecipients(client *smtp.Client, to []string) error {
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return err
		}
	}
	return nil
}

func writeAndClose(client *smtp.Client, message []byte) error {
	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write(message)
	if err != nil {
		return err
	}

	err = writer.Close()
	return err
}

func SendEmail(client *smtp.Client, email *EmailMessage) error {
	err := setMail(client, email.from)
	if err != nil {
		return err
	}

	err = setRecipients(client, email.to)
	if err != nil {
		return err
	}

	err = writeAndClose(client, email.Prepare())
	if err != nil {
		return err
	}

	err = client.Quit()
	return err
}

func (es *EmailSenderServiceImpl) SendExchangeRate(
	exchangeRate float32,
	emailAddresses []string,
) int {
	var config Config
	err := config.loadFromYamlFile(es.ConfigFilePath)
	if err != nil {
		log.Fatal(err)
	}

	client := NewSMTPClient(config.SMTP)
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
