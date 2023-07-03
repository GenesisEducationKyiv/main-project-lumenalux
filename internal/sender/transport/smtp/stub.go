package smtp

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
)

type StubWriteCloser struct {
	Writer io.Writer
	Closer io.Closer
}

func (m *StubWriteCloser) Close() error {
	return nil
}

func (m *StubWriteCloser) Write(data []byte) (int, error) {
	return 0, nil
}

type StubSMTPClient struct {
	authCalled bool
	quitCalled bool
	dataCalled bool
	mailCalled bool
	rcptCalled bool

	authErr error
	quitErr error
	dataErr error
	MailErr error
	rcptErr error

	writer io.WriteCloser
}

func (m *StubSMTPClient) Auth(a smtp.Auth) error {
	m.authCalled = true
	return m.authErr
}

func (m *StubSMTPClient) Quit() error {
	m.quitCalled = true
	return m.quitErr
}

func (m *StubSMTPClient) Data() (io.WriteCloser, error) {
	m.dataCalled = true
	m.writer = &StubWriteCloser{}
	return m.writer, m.dataErr
}

func (m *StubSMTPClient) Mail(from string) error {
	m.mailCalled = true
	return m.MailErr
}

func (m *StubSMTPClient) Rcpt(to string) error {
	m.rcptCalled = true
	return m.rcptErr
}

type StubDialer struct {
	Err error
}

func (d *StubDialer) Dial(network string, addr string, config *tls.Config) (*tls.Conn, error) {
	return nil, d.Err
}

type StubSMTPClientFactory struct {
	Client *StubSMTPClient
	Err    error
}

func (f StubSMTPClientFactory) NewClient(
	conn net.Conn,
	host string,
) (SMTPConnectionClient, error) {
	return f.Client, f.Err
}
