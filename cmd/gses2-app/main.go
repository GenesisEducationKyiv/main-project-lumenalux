package main

import (
	"context"
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
	"gses2-app/pkg/repository/userrepo"
	"gses2-app/pkg/storage"
)

func main() {
	ctx := context.Background()
	config, err := config.Load(ctx)
	if err != nil {
		log.Printf("Error, config wasn't loaded: %s", err)
		os.Exit(0)
	}

	senderService, err := createSenderService(&config)
	if err != nil {
		log.Printf("Connection error: %s", err)
		os.Exit(0)
	}

	rateService := createRateService(&config)
	subscriptionService := createSubscriptionService(&config)

	appController := controller.NewAppController(
		rateService,
		subscriptionService,
		senderService,
	)

	mux := registerRoutes(appController)
	startServer(config.HTTP.Port, mux)
}

func createRateService(config *config.Config) *rate.Service {
	httpClient := &http.Client{Timeout: config.HTTP.Timeout}
	exchangeRateProvider := kuna.NewKunaProvider(
		config.KunaAPI, httpClient,
	)
	return rate.NewService(exchangeRateProvider)
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

func createSubscriptionService(config *config.Config) *subscription.Service {
	storageCSV := storage.NewCSVStorage(config.Storage.Path)
	userRepository := userrepo.NewUserRepository(storageCSV)

	return subscription.NewService(userRepository)
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
