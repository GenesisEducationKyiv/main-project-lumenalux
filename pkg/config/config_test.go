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
		yamlName       string
		yamlContent    string
		expectedConfig *Config
		expectedErr    error
	}{
		{
			name:     "Load SMTP config",
			yamlName: "test_config.yaml",
			yamlContent: `
smtp:
  host: smtp.example.com
  port: 587
  user: user@example.com
  password: secret
`,
			expectedConfig: &Config{
				SMTP: SMTPConfig{
					Host:     "smtp.example.com",
					Port:     587,
					User:     "user@example.com",
					Password: "secret",
				},
			},
		},
		{
			name:     "Load email config",
			yamlName: "test_config.yaml",
			yamlContent: `
email:
  from: user@example.com
  subject: Test Subject
  body: Test Body
`,
			expectedConfig: &Config{
				Email: EmailConfig{
					From:    "user@example.com",
					Subject: "Test Subject",
					Body:    "Test Body",
				},
			},
		},
		{
			name:     "Load storage config",
			yamlName: "test_config.yaml",
			yamlContent: `
storage:
  path: /var/data
`,
			expectedConfig: &Config{
				Storage: StorageConfig{
					Path: "/var/data",
				},
			},
		},
		{
			name:     "Load HTTP config",
			yamlName: "test_config.yaml",
			yamlContent: `
http:
  port: "8080"
  timeout_in_seconds: 10
`,
			expectedConfig: &Config{
				HTTP: HTTPConfig{
					Port:    "8080",
					Timeout: 10,
				},
			},
		},
		{
			name:     "Load Kuna API config",
			yamlName: "test_config.yaml",
			yamlContent: `
kuna_api:
  url: https://api.example.com
`,
			expectedConfig: &Config{
				KunaAPI: KunaAPIConfig{
					URL: "https://api.example.com",
				},
			},
		},
		{
			name:        "Non existent file",
			yamlName:    "nonexistentfile.yaml",
			expectedErr: ErrReadFile,
		},
		{
			name:        "Invalid YAML",
			yamlName:    "config_test",
			yamlContent: `:\C`,
			expectedErr: ErrUnmarshalYAML,
		},
	}

	for _, td := range testData {
		t.Run(td.name, func(t *testing.T) {
			runTest(t, td)
		})
	}
}

func runTest(t *testing.T, td struct {
	name           string
	yamlName       string
	yamlContent    string
	expectedConfig *Config
	expectedErr    error
}) {
	var fileName string

	if td.expectedErr != nil && errors.Is(td.expectedErr, ErrReadFile) {
		fileName = td.yamlName
	} else {
		fileName = prepareTempFile(t, td.yamlName, td.yamlContent)
		defer os.Remove(fileName)
	}

	config, err := Load(fileName)
	compareResults(t, &config, err, td.expectedConfig, td.expectedErr)
}

func prepareTempFile(t *testing.T, yamlName string, yamlContent string) string {
	tempFile, err := os.CreateTemp("", yamlName)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	_, err = io.WriteString(tempFile, yamlContent)
	if err != nil {
		t.Fatalf("Failed to write test data to the temporary file: %v", err)
	}

	return tempFile.Name()
}

func compareResults(
	t *testing.T,
	config *Config,
	err error,
	expectedConfig *Config,
	expectedErr error,
) {
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Expected error %v but got %v", expectedErr, err)
	}

	if expectedConfig != nil && *config != *expectedConfig {
		t.Errorf("Unexpected configuration. Got: %+v, Expected: %+v", *config, *expectedConfig)
	}
}
