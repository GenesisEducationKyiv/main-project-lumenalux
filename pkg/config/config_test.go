package config

import (
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
			name: "ValidConfig",
			yaml: `
smtp:
  host: smtp.example.com
  port: 587
  user: user@example.com
  password: secret

email:
  from: user@example.com
  subject: Test Subject
  body: Test Body

storage:
  path: /var/data

http:
  port: "8080"
  timeout_in_seconds: 10

kuna_api:
  url: https://api.example.com
`,
			expectedConfig: Config{
				SMTP: SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					User:     "user@example.com",
					Password: "secret",
				},
				Email: EmailConfig{
					From:    "user@example.com",
					Subject: "Test Subject",
					Body:    "Test Body",
				},
				Storage: StorageConfig{
					Path: "/var/data",
				},
				HTTP: HTTPConfig{
					Port:    "8080",
					Timeout: 10,
				},
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
