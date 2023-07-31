package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gses2-app/internal/controller"
	"gses2-app/internal/rate"
	"gses2-app/internal/rate/provider/kuna"
	"gses2-app/internal/sender"
	"gses2-app/internal/sender/provider/email"
	"gses2-app/internal/sender/transport/smtp"
	"gses2-app/internal/subscription"
	"gses2-app/internal/transport"
	"gses2-app/pkg/config"
	"gses2-app/pkg/storage"
)

const _configPrefix = "GSES2_APP"

func main() {
	config, err := config.Load(_configPrefix)
	if err != nil {
		log.Printf("Error, config wasn't loaded: %s", err)
		os.Exit(0)
	}

	rateService, subscriptionService, senderService, err := createServices(&config)
	appController := controller.NewAppController(
		rateService,
		subscriptionService,
		senderService,
	)

	if err != nil {
		log.Printf("Connection error: %s", err)
		os.Exit(0)
	}

	mux := registerRoutes(appController)
	startServer(config.HTTP.Port, mux)
}

func createServices(config *config.Config) (
	*rate.Service,
	*subscription.Service,
	*sender.Service,
	error,
) {
	senderService, err := createSenderService(config)
	if err != nil {
		return nil, nil, nil, err
	}

	httpClient := &http.Client{Timeout: config.HTTP.Timeout}
	exchangeRateProvider := kuna.NewKunaProvider(
		config.KunaAPI, httpClient,
	)
	rateService := rate.NewService(exchangeRateProvider)

	emailSubscriptionService := subscription.NewService(storage.NewCSVStorage(config.Storage.Path))

	return rateService, emailSubscriptionService, senderService, nil
}

func createSenderService(
	config *config.Config,
) (*sender.Service, error) {
	emailSenderProvider, err := email.NewProvider(
		config,
		&smtp.TLSConnectionDialerImpl{},
		&smtp.SMTPClientFactoryImpl{},
	)

	if err != nil {
		return nil, err
	}

	return sender.NewService(emailSenderProvider), nil
}

func registerRoutes(appController *controller.AppController) *http.ServeMux {
	router := transport.NewHTTPRouter(appController)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux)

	return mux
}

func startServer(port string, handler http.Handler) {
	log.Printf("Starting server on port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Fatal(err)
	}
}
