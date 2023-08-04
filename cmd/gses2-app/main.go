package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"gses2-app/internal/core/port"
	"gses2-app/internal/core/service/rate"
	"gses2-app/internal/core/service/sender"
	"gses2-app/internal/core/service/subscription"
	"gses2-app/internal/handler/httpcontroller"
	"gses2-app/internal/handler/router"
	"gses2-app/internal/repository/config"
	"gses2-app/internal/repository/logger/rabbit"
	"gses2-app/internal/repository/rate/rest/binance"
	"gses2-app/internal/repository/rate/rest/coingecko"
	"gses2-app/internal/repository/rate/rest/kuna"
	"gses2-app/internal/repository/sender/email"
	"gses2-app/internal/repository/sender/smtp"
	"gses2-app/internal/repository/storage"
)

const (
	_configPrefix     = "GSES2_APP"
	_rabbitMQQueueUrl = "amqp://guest:guest@amqp/"
)

func main() {
	ctx := context.Background()
	conn, ch, q, err := rabbit.ConnectToRabbitMQ(_rabbitMQQueueUrl)
	if err != nil {
		log.Printf("Error, cannot connect to RabbitMQ: %s", err)
		os.Exit(1)
	}

	defer conn.Close()
	defer ch.Close()

	logger := rabbit.NewLogger(ctx, ch, q)

	consumer, err := rabbit.NewConsumer(ch, q)
	if err != nil {
		log.Printf("Error, cannot create logger consumer: %s", err)
		os.Exit(1)
	}

	loging := make(chan bool)

	go consumer()

	config, err := config.Load(_configPrefix)
	if err != nil {
		logger.Errorf("Error, config wasn't loaded: %s", err)
		os.Exit(1)
	}

	senderService, err := createSenderService(&config)
	if err != nil {
		logger.Errorf("Connection error: %s", err)
		os.Exit(1)
	}

	rateService := createRateService(&config)
	subscriptionService := createSubscriptionService(&config)

	appController := httpcontroller.NewAppController(
		rateService,
		subscriptionService,
		senderService,
	)

	mux := registerRoutes(appController)
	startServer(logger, config.HTTP.Port, mux)

	<-loging
}

func createRateService(config *config.Config) *rate.Service {

	httpClient := &http.Client{Timeout: config.HTTP.Timeout}

	BinanceRateProvider := binance.NewProvider(
		config.BinanceAPI, httpClient,
	)

	KunaRateProvider := kuna.NewProvider(
		config.KunaAPI, httpClient,
	)

	CoingeckoRateProvider := coingecko.NewProvider(
		config.CoingeckoAPI, httpClient,
	)

	return rate.NewService(
		BinanceRateProvider,
		CoingeckoRateProvider,
		KunaRateProvider,
	)
}

func createSenderService(
	config *config.Config,
) (*sender.Service, error) {
	emailSenderProvider, err := email.NewProvider(
		&email.EmailSenderConfig{
			SMTP:  config.SMTP,
			Email: config.Email,
		},
		&smtp.TLSConnectionDialerImpl{},
		&smtp.SMTPClientFactoryImpl{},
	)

	if err != nil {
		return nil, err
	}

	return sender.NewService(emailSenderProvider), nil
}

func createSubscriptionService(config *config.Config) *subscription.Service {
	storageCSV := storage.NewCSVStorage(config.Storage.Path)
	userRepository := port.NewUserRepository(storageCSV)

	return subscription.NewService(userRepository)
}

func registerRoutes(appController *httpcontroller.AppController) *http.ServeMux {
	router := router.NewHTTPRouter(appController)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux)

	return mux
}

func startServer(logger port.Logger, port string, handler http.Handler) {
	logger.Infof("Starting server on port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}
