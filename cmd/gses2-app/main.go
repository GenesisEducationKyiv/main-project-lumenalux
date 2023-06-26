package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gses2-app/internal/controllers"
	"gses2-app/internal/email"
	"gses2-app/internal/rate"
	"gses2-app/internal/subscription"
	"gses2-app/internal/transport"
	"gses2-app/pkg/config"
	"gses2-app/pkg/storage"
)

func main() {
	config, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatalf("Error, config wasn't loaded: %s", err)
	}

	rateService, emailSubscriptionService, emailSenderService, err := createServices(&config)
	controller := controllers.NewAppController(
		rateService,
		emailSubscriptionService,
		emailSenderService,
	)

	if err != nil {
		log.Fatalf("Connection error: %s", err)
	}

	mux := registerRoutes(controller)
	startServer(config.HTTP.Port, mux)
}

func createServices(config *config.Config) (
	*rate.Service,
	*subscription.Service,
	*email.SenderService,
	error,
) {
	httpClient := &http.Client{Timeout: config.HTTP.Timeout * time.Second}

	rateService := rate.NewService(rate.NewKunaProvider(config.KunaAPI, httpClient))

	emailSubscriptionService := subscription.NewService(storage.NewCSVStorage(config.Storage.Path))

	emailSenderService, err := email.NewSenderService(
		config,
		&email.TLSConnectionDialerImpl{},
		&email.SMTPClientFactoryImpl{},
	)

	if err != nil {
		return nil, nil, nil, err
	}

	return rateService, emailSubscriptionService, emailSenderService, err
}

func registerRoutes(controller *controllers.AppController) *http.ServeMux {
	router := transport.NewHTTPRouter(controller)

	mux := http.NewServeMux()
	router.RegisterRoutes(mux)

	return mux
}

func startServer(port string, handler http.Handler) {
	fmt.Printf("Starting server on port %s\n", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Fatal(err)
	}
}
