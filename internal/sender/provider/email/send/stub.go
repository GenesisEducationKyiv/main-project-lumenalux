package send

import (
	"io"
)

type StubSMTPClient struct {
	fromCalledWith    string
	rcptCalledWith    []string
	rcptShouldReturn  error
	dataCalled        bool
	quitCalled        bool
	writeCalledWith   []byte
	writeShouldReturn error
	mailShouldReturn  error
}

func (m *StubSMTPClient) Mail(from string) error {
	m.fromCalledWith = from
	return m.mailShouldReturn
}

func (m *StubSMTPClient) Rcpt(to string) error {
	m.rcptCalledWith = append(m.rcptCalledWith, to)
	return m.rcptShouldReturn
}

func (m *StubSMTPClient) Data() (wc io.WriteCloser, err error) {
	m.dataCalled = true
	if m.writeShouldReturn != nil {
		return nil, m.writeShouldReturn
	}
	return m, nil
}

func (m *StubSMTPClient) Write(p []byte) (n int, err error) {
	m.writeCalledWith = p
	return len(p), nil
}

func (m *StubSMTPClient) Close() error {
	return nil
}

func (m *StubSMTPClient) Quit() error {
	m.quitCalled = true
	return nil
}
