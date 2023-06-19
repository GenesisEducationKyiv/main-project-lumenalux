package email

import (
	"crypto/tls"
	"net/smtp"
	"strconv"

	"gses2-app/pkg/config"
)

type SMTPClient struct {
	host     string
	port     int
	user     string
	password string
}

func NewSMTPClient(config config.SMTPConfig) *SMTPClient {
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
