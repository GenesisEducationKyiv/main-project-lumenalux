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
	From    string `env:"FROM,default=no.reply@currency.info.api"`
	Subject string `env:"SUBJECT,default=BTC to UAH exchange rate"`
	Body    string `env:"BODY,default=The BTC to UAH exchange rate is {{.Rate}} UAH per BTC"`
}

type SMTPConfig struct {
	Host     string `env:"HOST,required"`
	Port     int    `env:"PORT,default=465"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
}

type StorageConfig struct {
	Path string `env:"PATH,default=./storage/storage.csv"`
}

type HTTPConfig struct {
	Port    string        `env:"PORT,default=8080"`
	Timeout time.Duration `env:"TIMEOUT,default=10s"`
}

type KunaAPIConfig struct {
	URL         string     `env:"URL,default=https://api.kuna.io/v3/tickers?symbols=btcuah"`
	DefaultRate types.Rate `env:"DEFAULT_RATE,default=0"`
}
