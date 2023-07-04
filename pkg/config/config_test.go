package config

import (
	"context"
	"errors"
	"gses2-app/pkg/types"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

var (
	_defaultEnvVariables = map[string]string{
		"GSES2_APP_SMTP_HOST":             "www.default.com",
		"GSES2_APP_SMTP_USER":             "default@user.com",
		"GSES2_APP_SMTP_PASSWORD":         "defaultpassword",
		"GSES2_APP_SMTP_PORT":             "465",
		"GSES2_APP_EMAIL_FROM":            "no.reply@test.info.api",
		"GSES2_APP_EMAIL_SUBJECT":         "BTC to UAH exchange rate",
		"GSES2_APP_EMAIL_BODY":            "The BTC to UAH rate is {{.Rate}}",
		"GSES2_APP_STORAGE_PATH":          "./storage/storage.csv",
		"GSES2_APP_HTTP_PORT":             "8080",
		"GSES2_APP_HTTP_TIMEOUT":          "10s",
		"GSES2_APP_KUNA_API_URL":          "https://www.example.com",
		"GSES2_APP_KUNA_API_DEFAULT_RATE": "0",
	}
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		envVars        map[string]string
		updateExpected func(t *testing.T, c Config) Config
		expectedErr    error
	}{
		{
			name: "All required variables provided",
			envVars: map[string]string{
				"GSES2_APP_SMTP_HOST":     "smtp.example.com",
				"GSES2_APP_SMTP_USER":     "user@example.com",
				"GSES2_APP_SMTP_PASSWORD": "secret",
			},
			updateExpected: func(t *testing.T, c Config) Config {
				c.SMTP.Host = "smtp.example.com"
				c.SMTP.User = "user@example.com"
				c.SMTP.Password = "secret"
				return c
			},
		},
		{
			name:        "Missing required variables",
			envVars:     map[string]string{},
			expectedErr: ErrLoadEnvVariable,
		},
		{
			name: "Override default variable",
			envVars: initEnvVariables(map[string]string{
				"GSES2_APP_EMAIL_FROM": "override@example.com",
			}),
			updateExpected: func(t *testing.T, c Config) Config {
				c = addDefaultConfigVariables(t, c)
				c.Email.From = "override@example.com"
				return c
			},
		},
		{
			name: "Override multiple default variables",
			envVars: initEnvVariables(map[string]string{
				"GSES2_APP_EMAIL_FROM":   "override@example.com",
				"GSES2_APP_SMTP_PORT":    "999",
				"GSES2_APP_STORAGE_PATH": "/new/path",
				"GSES2_APP_HTTP_TIMEOUT": "15s",
				"GSES2_APP_KUNA_API_URL": "https://new.api.url",
			}),
			updateExpected: func(t *testing.T, c Config) Config {
				c = addDefaultConfigVariables(t, c)
				c.Email.From = "override@example.com"
				c.SMTP.Port = 999
				c.Storage.Path = "/new/path"
				c.HTTP.Timeout = 15 * time.Second
				c.KunaAPI.URL = "https://new.api.url"
				return c
			},
		},
		{
			name: "Missing one required variable",
			envVars: map[string]string{
				"GSES2_APP_SMTP_HOST": "smtp.example.com",
				"GSES2_APP_SMTP_USER": "user@example.com",
			},
			expectedErr: ErrLoadEnvVariable,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			initTestEnvironment(t, tt.envVars)

			ctx := context.Background()
			config, err := Load(ctx)

			if tt.expectedErr != nil {
				log.Println(err, tt.expectedErr)
				if !errors.Is(err, tt.expectedErr) {
					t.Fatalf("In test %v\nExpected:\n%v\nbut got:\n%v\n", t.Name(), tt.expectedErr, err)
				}
				require.ErrorIs(t, err, tt.expectedErr)
				return
			}

			require.NoError(t, err)
			expectedConfig := tt.updateExpected(t, defaultConfig())
			require.Equal(t, expectedConfig, config)
		})
	}
}

func initTestEnvironment(t *testing.T, envVars map[string]string) {
	for key := range _defaultEnvVariables {
		t.Setenv(key, "")
		os.Unsetenv(key)
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}
}

func defaultConfig() Config {
	return Config{
		SMTP: SMTPConfig{
			Port: 465,
		},
		Email: EmailConfig{
			From:    "no.reply@currency.info.api",
			Subject: "BTC to UAH exchange rate",
			Body:    "The BTC to UAH exchange rate is {{.Rate}} UAH per BTC",
		},
		Storage: StorageConfig{
			Path: "./storage/storage.csv",
		},
		HTTP: HTTPConfig{
			Port:    "8080",
			Timeout: 10 * time.Second,
		},
		KunaAPI: KunaAPIConfig{
			URL:         "https://api.kuna.io/v3/tickers?symbols=btcuah",
			DefaultRate: 0,
		},
	}
}

func addDefaultConfigVariables(t *testing.T, c Config) Config {
	c.SMTP.Host = _defaultEnvVariables["GSES2_APP_SMTP_HOST"]
	c.SMTP.User = _defaultEnvVariables["GSES2_APP_SMTP_USER"]
	c.SMTP.Password = _defaultEnvVariables["GSES2_APP_SMTP_PASSWORD"]
	c.SMTP.Port = parseSMTPPort(t, _defaultEnvVariables["GSES2_APP_SMTP_PORT"])
	c.Email.From = _defaultEnvVariables["GSES2_APP_EMAIL_FROM"]
	c.Email.Subject = _defaultEnvVariables["GSES2_APP_EMAIL_SUBJECT"]
	c.Email.Body = _defaultEnvVariables["GSES2_APP_EMAIL_BODY"]
	c.Storage.Path = _defaultEnvVariables["GSES2_APP_STORAGE_PATH"]
	c.HTTP.Port = _defaultEnvVariables["GSES2_APP_HTTP_PORT"]
	c.HTTP.Timeout, _ = time.ParseDuration(_defaultEnvVariables["GSES2_APP_HTTP_TIMEOUT"])
	c.KunaAPI.URL = _defaultEnvVariables["GSES2_APP_KUNA_API_URL"]
	c.KunaAPI.DefaultRate = parseKunaAPIDefaultRate(
		t,
		_defaultEnvVariables["GSES2_APP_KUNA_API_DEFAULT_RATE"],
	)

	return c
}

func parseSMTPPort(t *testing.T, strPort string) int {
	SMTPPort, err := strconv.Atoi(strPort)
	if err != nil {
		t.Fatal("cannot convert default SMTP port value")
	}

	return SMTPPort
}

func parseKunaAPIDefaultRate(t *testing.T, strRate string) types.Rate {
	rate, err := strconv.ParseFloat(strRate, 32)
	if err != nil {
		t.Fatal("cannot convert default Kuna API rate")
	}

	return types.Rate(rate)
}

func initEnvVariables(newEnvVariables map[string]string) map[string]string {
	envVariables := map[string]string{}
	maps.Copy(envVariables, _defaultEnvVariables)
	for k, v := range newEnvVariables {
		envVariables[k] = v
	}

	return envVariables
}
