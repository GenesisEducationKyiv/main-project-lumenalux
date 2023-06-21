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
	"gses2-app/pkg/config"
	"gses2-app/pkg/storage"
)

func main() {
	err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	config := config.Current()

	httpClient := &http.Client{Timeout: time.Second * 10}
	exchangeRateService := rate.NewService(
		rate.NewKunaProvider(httpClient),
	)

	emailSubscriptionService := subscription.NewService(
		storage.NewCSVStorage(config.Storage.Path),
	)

	emailSenderService := email.NewSenderService(
		&email.TLSConnectionDialerImpl{},
		&email.SMTPClientFactoryImpl{},
	)

	controller := controllers.NewAppController(
		exchangeRateService,
		emailSubscriptionService,
		emailSenderService,
	)

	http.HandleFunc("/api/rate", controller.GetRate)
	http.HandleFunc("/api/subscribe", controller.SubscribeEmail)
	http.HandleFunc("/api/sendEmails", controller.SendEmails)

	message := fmt.Sprintf("Starting server on port %s", config.HTTP.Port)
	fmt.Println(message)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", config.HTTP.Port), nil))
}
