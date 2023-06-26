package config

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	testData := []struct {
		name           string
		yaml           string
		expectedConfig Config
	}{
		{
			name: "Load SMTP config",
			yaml: `
smtp:
  host: smtp.example.com
  port: 587
  user: user@example.com
  password: secret
`,
			expectedConfig: Config{
				SMTP: SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					User:     "user@example.com",
					Password: "secret",
				},
			},
		},
		{
			name: "Load email config",
			yaml: `
email:
  from: user@example.com
  subject: Test Subject
  body: Test Body
`,
			expectedConfig: Config{
				Email: EmailConfig{
					From:    "user@example.com",
					Subject: "Test Subject",
					Body:    "Test Body",
				},
			},
		},
		{
			name: "Load storage config",
			yaml: `
storage:
  path: /var/data
`,
			expectedConfig: Config{
				Storage: StorageConfig{
					Path: "/var/data",
				},
			},
		},
		{
			name: "Load HTTP config",
			yaml: `
http:
  port: "8080"
  timeout_in_seconds: 10
`,
			expectedConfig: Config{
				HTTP: HTTPConfig{
					Port:    "8080",
					Timeout: 10,
				},
			},
		},
		{
			name: "Load Kuna API config",
			yaml: `
kuna_api:
  url: https://api.example.com
`,
			expectedConfig: Config{
				KunaAPI: KunaAPIConfig{
					URL: "https://api.example.com",
				},
			},
		},
	}

	for _, testData := range testData {
		t.Run(testData.name, func(t *testing.T) {

			tempFile, err := os.CreateTemp("", "config_test")
			if err != nil {
				t.Fatalf("Failed to create temporary file: %v", err)
			}
			defer os.Remove(tempFile.Name())

			_, err = io.WriteString(tempFile, testData.yaml)
			if err != nil {
				t.Fatalf("Failed to write test data to the temporary file: %v", err)
			}

			config, err := Load(tempFile.Name())
			if err != nil {
				t.Fatalf("Failed to load configuration: %v", err)
			}

			if config != testData.expectedConfig {
				t.Errorf("Unexpected configuration. Got: %+v, Expected: %+v", config, testData.expectedConfig)
			}
		})
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	_, err := Load("nonexistentfile.yaml")
	if err == nil || !errors.Is(err, ErrReadFile) {
		t.Fatalf("Expected error %v but got %v", ErrReadFile, err)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	tempFile, err := os.CreateTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = io.WriteString(tempFile, `:\C`)
	if err != nil {
		t.Fatalf("Failed to write test data to the temporary file: %v", err)
	}

	_, err = Load(tempFile.Name())
	if err == nil || !errors.Is(err, ErrUnmarshalYAML) {
		t.Fatalf("Expected error %v but got %v", ErrUnmarshalYAML, err)
	}
}
