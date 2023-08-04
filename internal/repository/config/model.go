package config

import (
	"gses2-app/internal/handler/router"
	"gses2-app/internal/repository/rate/rest/binance"
	"gses2-app/internal/repository/rate/rest/coingecko"
	"gses2-app/internal/repository/rate/rest/kuna"
	"gses2-app/internal/repository/sender/email/send"
	"gses2-app/internal/repository/sender/smtp"
	"gses2-app/internal/repository/storage"
)

type Config struct {
	SMTP         smtp.SMTPConfig
	Email        send.EmailConfig
	Storage      storage.StorageConfig
	HTTP         router.HTTPConfig
	KunaAPI      kuna.KunaAPIConfig
	BinanceAPI   binance.BinanceAPIConfig
	CoingeckoAPI coingecko.CoingeckoAPIConfig
}
