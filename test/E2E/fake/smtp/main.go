package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"

	"github.com/mhale/smtpd"
)

const (
	AppName = "SMTP server"

	Hostname = "localhost"
	Port     = ":1025"
	Username = "test"
	Password = "password"

	OneMBSize = 1024 * 1024
)

func authHandler(remoteAddr net.Addr, mechanism string, username []byte, password []byte, shared []byte) (bool, error) {
	return string(username) == Username && string(password) == Password, nil
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
		log.Fatalf("Failed to load key pair: %v", err)
	}

	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ServerName:         Hostname,
		InsecureSkipVerify: true,
	}
}

func main() {
	server := &smtpd.Server{
		Addr:         Port,
		Hostname:     Hostname,
		Handler:      mailHandler,
		AuthHandler:  authHandler,
		Appname:      AppName,
		MaxSize:      OneMBSize,
		AuthRequired: true,
		TLSConfig:    getTLSConfig(),
	}

	fmt.Printf("Starting server on %s\n", Port)

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}

	tlsListener := tls.NewListener(listener, server.TLSConfig)
	if err := server.Serve(tlsListener); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
