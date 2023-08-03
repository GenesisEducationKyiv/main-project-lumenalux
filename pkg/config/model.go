package config

import (
	"gses2-app/internal/rate/provider/binance"
	"gses2-app/internal/rate/provider/coingecko"
	"gses2-app/internal/rate/provider/kuna"
	"gses2-app/internal/sender/provider/email/message"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/pkg/storage"
	"time"
)

type Config struct {
	SMTP         smtp.SMTPConfig
	Email        message.EmailConfig
	Storage      storage.StorageConfig
	HTTP         HTTPConfig
	KunaAPI      kuna.KunaAPIConfig
	BinanceAPI   binance.BinanceAPIConfig
	CoingeckoAPI coingecko.CoingeckoAPIConfig
}

type HTTPConfig struct {
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}
