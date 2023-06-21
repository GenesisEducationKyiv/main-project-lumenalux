package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/smtp"
	"strconv"

	"gses2-app/pkg/config"
)

type TLSConnectionDialer interface {
	Dial(network, addr string, config *tls.Config) (*tls.Conn, error)
}

type TLSConnectionDialerImpl struct{}

func (d TLSConnectionDialerImpl) Dial(network, addr string, config *tls.Config) (conn *tls.Conn, err error) {
	return tls.Dial(network, addr, config)
}

type SMTPConnectionClient interface {
	Auth(a smtp.Auth) error
	Quit() error
	Data() (io.WriteCloser, error)
	Mail(string) error
	Rcpt(string) error
}

type SMTPClientFactory interface {
	NewClient(conn net.Conn, host string) (SMTPConnectionClient, error)
}

type SMTPClientFactoryImpl struct{}

func (f SMTPClientFactoryImpl) NewClient(
	conn net.Conn,
	host string,
) (SMTPConnectionClient, error) {
	return smtp.NewClient(conn, host)
}

type SMTPClient struct {
	host              string
	port              int
	user              string
	password          string
	dialer            TLSConnectionDialer
	smtpClientFactory SMTPClientFactory
}

func NewSMTPClient(
	config config.SMTPConfig,
	dialer TLSConnectionDialer,
	factory SMTPClientFactory,
) *SMTPClient {
	return &SMTPClient{
		host:              config.Host,
		port:              config.Port,
		user:              config.User,
		password:          config.Password,
		dialer:            dialer,
		smtpClientFactory: factory,
	}
}

func (c *SMTPClient) createTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         c.host,
	}
}

func (c *SMTPClient) createConnection(tlsConfig *tls.Config) (*tls.Conn, error) {
	conn, err := c.dialer.Dial(
		"tcp",
		fmt.Sprintf("%s:%s", c.host, strconv.Itoa(c.port)),
		tlsConfig,
	)
	return conn, err
}

func (c *SMTPClient) createSMTPClient(conn *tls.Conn) (SMTPConnectionClient, error) {
	client, err := c.smtpClientFactory.NewClient(conn, c.host)
	return client, err
}

func (c *SMTPClient) authenticate(client SMTPConnectionClient) error {
	auth := smtp.PlainAuth("", c.user, c.password, c.host)
	return client.Auth(auth)
}

func (c *SMTPClient) Connect() (SMTPConnectionClient, error) {
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
