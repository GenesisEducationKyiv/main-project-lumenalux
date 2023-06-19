package config

import (
	"os"
	"reflect"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestLoad(t *testing.T) {
	sampleConfig := Config{
		Email: EmailConfig{
			From:    "test@example.com",
			Subject: "test subject",
			Body:    "test body",
		},
		SMTP: SMTPConfig{
			Host:     "smtp.example.com",
			Port:     587,
			User:     "user",
			Password: "password",
		},
		Storage: StorageConfig{
			Path: "/tmp",
		},
		HTTP: HTTPConfig{
			Port: "8080",
		},
		KunaAPI: KunaAPIConfig{
			URL: "https://kuna.io/api/v2",
		},
	}

	data, err := yaml.Marshal(&sampleConfig)
	if err != nil {
		t.Fatalf("failed to marshal sample config: %v", err)
	}

	tmpfile, err := os.CreateTemp("", "testconfig")
	if err != nil {
		t.Fatalf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatalf("failed to write to temporary file: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("failed to close temporary file: %v", err)
	}

	if err := Load(tmpfile.Name()); err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if !reflect.DeepEqual(Current(), sampleConfig) {
		t.Errorf("loaded config does not match expected config")
	}
}
