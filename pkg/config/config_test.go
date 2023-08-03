package config

import (
	"errors"
	"gses2-app/internal/rate/provider/binance"
	"gses2-app/internal/rate/provider/coingecko"
	"gses2-app/internal/rate/provider/kuna"
	"gses2-app/internal/sender/provider/email/message"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/pkg/storage"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

const _configPrefix = "GSES2_APP"

var (
	_defaultEnvVariables = map[string]string{
		"GSES2_APP_SMTP_HOST":        "www.default.com",
		"GSES2_APP_SMTP_USER":        "default@user.com",
		"GSES2_APP_SMTP_PASSWORD":    "defaultpassword",
		"GSES2_APP_SMTP_PORT":        "465",
		"GSES2_APP_EMAIL_FROM":       "no.reply@test.info.api",
		"GSES2_APP_EMAIL_SUBJECT":    "BTC to UAH exchange rate",
		"GSES2_APP_EMAIL_BODY":       "The BTC to UAH rate is {{.Rate}}",
		"GSES2_APP_STORAGE_PATH":     "./storage/storage.csv",
		"GSES2_APP_HTTP_PORT":        "8080",
		"GSES2_APP_HTTP_TIMEOUT":     "10s",
		"GSES2_APP_KUNAAPI_URL":      "https://www.example.com",
		"GSES2_APP_BINANCEAPI_URL":   "https://www.example.com",
		"GSES2_APP_COINGECKOAPI_URL": "https://www.example.com",
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
				"GSES2_APP_KUNAAPI_URL":  "https://new.api.url",
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
		t.Run(tt.name, func(t *testing.T) {
			initTestEnvironment(t, tt.envVars)

			config, err := Load(_configPrefix)

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
	// Clean up environment to allow each test
	// case to start with a clean environment
	for key := range _defaultEnvVariables {
		t.Setenv(key, "") // Set environment variable as null
		os.Unsetenv(key)  // Remove environment variable completely
	}

	for key, value := range envVars {
		t.Setenv(key, value)
	}
}

func defaultConfig() Config {
	return Config{
		SMTP: smtp.SMTPConfig{
			Port: 465,
		},
		Email: message.EmailConfig{
			From:    "no.reply@currency.info.api",
			Subject: "BTC to UAH exchange rate",
			Body:    "The BTC to UAH exchange rate is {{.Rate}} UAH per BTC",
		},
		Storage: storage.StorageConfig{
			Path: "./storage/storage.csv",
		},
		HTTP: HTTPConfig{
			Port:    "8080",
			Timeout: 10 * time.Second,
		},
		KunaAPI: kuna.KunaAPIConfig{
			URL: "https://api.kuna.io/v3/tickers?symbols=btcuah",
		},
		BinanceAPI: binance.BinanceAPIConfig{
			URL: "https://api.binance.com/api/v3/klines?symbol=BTCUAH&interval=1s&limit=1",
		},
		CoingeckoAPI: coingecko.CoingeckoAPIConfig{
			URL: "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah",
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
	c.HTTP.Timeout, _ = time.ParseDuration(
		_defaultEnvVariables["GSES2_APP_HTTP_TIMEOUT"],
	)

	c.BinanceAPI.URL = _defaultEnvVariables["GSES2_APP_BINANCEAPI_URL"]

	c.CoingeckoAPI.URL = _defaultEnvVariables["GSES2_APP_COINGECKOAPI_URL"]

	c.KunaAPI.URL = _defaultEnvVariables["GSES2_APP_KUNAAPI_URL"]

	return c
}

func parseSMTPPort(t *testing.T, strPort string) int {
	SMTPPort, err := strconv.Atoi(strPort)
	if err != nil {
		t.Fatal("cannot convert default SMTP port value")
	}

	return SMTPPort
}

func initEnvVariables(newEnvVariables map[string]string) map[string]string {
	envVariables := map[string]string{}
	maps.Copy(envVariables, _defaultEnvVariables)
	for k, v := range newEnvVariables {
		envVariables[k] = v
	}

	return envVariables
}
