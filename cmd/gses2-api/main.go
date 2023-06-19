package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gses2-app/controllers"
	"gses2-app/services"
)

func main() {
	httpClient := &http.Client{Timeout: time.Second * 10}
	exchangeRateService := services.NewExchangeRateService(httpClient)
	emailSubscriptionService := services.NewEmailSubscriptionService("./storage.csv")
	emailSenderService := services.NewEmailSenderService("./config.yaml")

	controller := controllers.NewAppController(
		exchangeRateService,
		&emailSubscriptionService,
		emailSenderService,
	)

	http.HandleFunc("/api/rate", controller.GetRate)
	http.HandleFunc("/api/subscribe", controller.SubscribeEmail)
	http.HandleFunc("/api/sendEmails", controller.SendEmails)

	fmt.Println("Starting server on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
