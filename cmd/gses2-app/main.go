package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gses2-app/internal/controller"
	"gses2-app/internal/rate"
	"gses2-app/internal/sender"
	"gses2-app/internal/subscription"
	"gses2-app/internal/transport"
	"gses2-app/pkg/config"
	"gses2-app/pkg/storage"
)

func main() {
	ctx := context.Background()
	config, err := config.Load(ctx)
	if err != nil {
		log.Printf("Error, config wasn't loaded: %s", err)
		os.Exit(0)
	}

	rateService, emailSubscriptionService, emailSenderService, err := createServices(&config)
	appController := controller.NewAppController(
		rateService,
		emailSubscriptionService,
		emailSenderService,
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
	httpClient := &http.Client{Timeout: config.HTTP.Timeout * time.Second}

	rateService := rate.NewService(rate.NewKunaProvider(config.KunaAPI, httpClient))

	emailSubscriptionService := subscription.NewService(storage.NewCSVStorage(config.Storage.Path))

	emailSenderService, err := sender.NewService(
		config,
		&sender.TLSConnectionDialerImpl{},
		&sender.SMTPClientFactoryImpl{},
	)

	if err != nil {
		return nil, nil, nil, err
	}

	return rateService, emailSubscriptionService, emailSenderService, err
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
