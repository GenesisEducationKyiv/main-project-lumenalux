package email

import (
	"crypto/tls"
	"io"
	"net"
	"net/smtp"
)

type MockSenderSMTPClient struct {
	fromCalledWith    string
	rcptCalledWith    []string
	dataCalled        bool
	quitCalled        bool
	writeCalledWith   []byte
	writeShouldReturn error
}

func (m *MockSenderSMTPClient) Mail(from string) error {
	m.fromCalledWith = from
	return nil
}

func (m *MockSenderSMTPClient) Rcpt(to string) error {
	m.rcptCalledWith = append(m.rcptCalledWith, to)
	return nil
}

func (m *MockSenderSMTPClient) Data() (wc io.WriteCloser, err error) {
	m.dataCalled = true
	return m, m.writeShouldReturn
}

func (m *MockSenderSMTPClient) Write(p []byte) (n int, err error) {
	m.writeCalledWith = p
	return len(p), m.writeShouldReturn
}

func (m *MockSenderSMTPClient) Close() error {
	return nil
}

func (m *MockSenderSMTPClient) Quit() error {
	m.quitCalled = true
	return nil
}

type MockSMTPClient struct {
	authCalled bool
	quitCalled bool
	dataCalled bool
	mailCalled bool
	rcptCalled bool

	authErr error
	quitErr error
	dataErr error
	mailErr error
	rcptErr error

	writer io.WriteCloser
}

func (m *MockSMTPClient) Auth(a smtp.Auth) error {
	m.authCalled = true
	return m.authErr
}

func (m *MockSMTPClient) Quit() error {
	m.quitCalled = true
	return m.quitErr
}

type MockWriteCloser struct {
	Writer io.Writer
	Closer io.Closer
}

func (m *MockWriteCloser) Close() error {
	return nil
}

func (m *MockWriteCloser) Write(data []byte) (int, error) {
	return 0, nil
}

func (m *MockSMTPClient) Data() (io.WriteCloser, error) {
	m.dataCalled = true
	m.writer = &MockWriteCloser{}
	return m.writer, m.dataErr
}

func (m *MockSMTPClient) Mail(from string) error {
	m.mailCalled = true
	return m.mailErr
}

func (m *MockSMTPClient) Rcpt(to string) error {
	m.rcptCalled = true
	return m.rcptErr
}

type MockDialer struct{}

func (d *MockDialer) Dial(network string, addr string, config *tls.Config) (*tls.Conn, error) {
	return nil, nil
}

type MockSMTPClientFactory struct {
	Client *MockSMTPClient
}

func (f MockSMTPClientFactory) NewClient(
	conn net.Conn,
	host string,
) (SMTPConnectionClient, error) {
	return f.Client, nil
}
