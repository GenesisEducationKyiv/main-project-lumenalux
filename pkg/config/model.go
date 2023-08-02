package config

import (
	"gses2-app/pkg/types"
	"time"
)

type Config struct {
	SMTP         SMTPConfig
	Email        EmailConfig
	Storage      StorageConfig
	HTTP         HTTPConfig
	KunaAPI      KunaAPIConfig
	BinanceAPI   BinanceAPIConfig
	CoingeckoAPI CoingeckoAPIConfig
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

type BinanceAPIConfig struct {
	URL         string     `default:"https://api.binance.com/api/v3/klines?symbol=BTCUAH&interval=1s&limit=1"`
	DefaultRate types.Rate `default:"0"`
}

type CoingeckoAPIConfig struct {
	URL         string     `default:"https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=uah"`
	DefaultRate types.Rate `default:"0"`
}
