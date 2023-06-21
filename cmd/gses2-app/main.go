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
	err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := config.Current()

	rateService, emailSubscriptionService, emailSenderService := createServices(&config)
	controller := controllers.NewAppController(
		rateService,
		emailSubscriptionService,
		emailSenderService,
	)

	mux := registerRoutes(controller)
	startServer(config.HTTP.Port, mux)
}

func createServices(config *config.Config) (rate.Service, subscription.Service, email.SenderService) {
	httpClient := &http.Client{Timeout: config.HTTP.Timeout * time.Second}

	rateService := rate.NewService(rate.NewKunaProvider(config.KunaAPI, httpClient))

	emailSubscriptionService := subscription.NewService(storage.NewCSVStorage(config.Storage.Path))

	emailSenderService := email.NewSenderService(
		&email.TLSConnectionDialerImpl{},
		&email.SMTPClientFactoryImpl{},
	)

	return rateService, emailSubscriptionService, emailSenderService
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
