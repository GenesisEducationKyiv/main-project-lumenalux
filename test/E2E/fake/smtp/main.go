package main

import (
	"bytes"
	"crypto/tls"
	"log"
	"net"
	"net/mail"
	"os"

	"github.com/mhale/smtpd"
)

const (
	_appName = "SMTP server"

	_hostname = "localhost"
	_port     = ":1025"
	_username = "test"
	_password = "password"

	_oneMBSize = 1024 * 1024
)

func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return string(username) == _username && string(password) == _password, nil
}

func mailHandler(origin net.Addr, from string, to []string, data []byte) error {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to parse mail: %v", err)
		return err
	}

	subject := msg.Header.Get("Subject")
	log.Printf("Received mail from %s for %s with subject %s", from, to[0], subject)
	return nil
}

func getTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Printf("Failed to load key pair: %v", err)
		os.Exit(0)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ServerName:         _hostname,
		InsecureSkipVerify: true,
	}
}

func main() {
	server := &smtpd.Server{
		Addr:         _port,
		Hostname:     _hostname,
		Handler:      mailHandler,
		AuthHandler:  authHandler,
		Appname:      _appName,
		MaxSize:      _oneMBSize,
		AuthRequired: true,
		TLSConfig:    getTLSConfig(),
	}

	log.Printf("Starting server on %s\n", _port)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Printf("Failed to create listener: %v", err)
		os.Exit(0)
	}

	tlsListener := tls.NewListener(listener, server.TLSConfig)
	if err := server.Serve(tlsListener); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
