package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gses2-app/internal/controllers"
	"gses2-app/internal/email"
	"gses2-app/internal/services"
	"gses2-app/pkg/config"
)

func main() {
	err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	httpClient := &http.Client{Timeout: time.Second * 10}
	exchangeRateService := services.NewExchangeRateService(httpClient)
	emailSubscriptionService := services.NewEmailSubscriptionService("./storage.csv")
	emailSenderService := email.NewSenderService()

	controller := controllers.NewAppController(
		exchangeRateService,
		&emailSubscriptionService,
		emailSenderService,
	)

	http.HandleFunc("/api/rate", controller.GetRate)
	http.HandleFunc("/api/subscribe", controller.SubscribeEmail)
	http.HandleFunc("/api/sendEmails", controller.SendEmails)

	fmt.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}