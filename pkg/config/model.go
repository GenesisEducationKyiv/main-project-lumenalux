package config

import (
	"gses2-app/pkg/types"
	"time"
)

type Config struct {
	SMTP    SMTPConfig    `env:",prefix=GSES2_APP_SMTP_"`
	Email   EmailConfig   `env:",prefix=GSES2_APP_EMAIL_"`
	Storage StorageConfig `env:",prefix=GSES2_APP_STORAGE_"`
	HTTP    HTTPConfig    `env:",prefix=GSES2_APP_HTTP_"`
	KunaAPI KunaAPIConfig `env:",prefix=GSES2_APP_KUNA_API_"`
}

type EmailConfig struct {
	From    string `default:"no.reply@currency.info.api"`
	Subject string `default:"BTC to UAH exchange rate"`
	Body    string `default:"The BTC to UAH exchange rate is {{.Rate}} UAH per BTC"`
}

type SMTPConfig struct {
	Host     string `required:"true"`
	Port     int    `default:"465"`
	User     string `required:"true"`
	Password string `required:"true"`
}

type StorageConfig struct {
	Path string `default:"./storage/storage.csv"`
}

type HTTPConfig struct {
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}

type KunaAPIConfig struct {
	URL         string     `default:"https://api.kuna.io/v3/tickers?symbols=btcuah"`
	DefaultRate types.Rate `default:"0"`
}
