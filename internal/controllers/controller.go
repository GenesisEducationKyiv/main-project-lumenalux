package controllers

import (
	"encoding/json"
	"net/http"
)

type SenderService interface {
	SendExchangeRate(float32, []string) error
}

type RateService interface {
	ExchangeRate() (float32, error)
}

type SubscriptionService interface {
	Subscribe(email string) error
	IsSubscribed(email string) (bool, error)
	Subscriptions() ([]string, error)
}

type AppController struct {
	ExchangeRateService      RateService
	EmailSubscriptionService SubscriptionService
	EmailSenderService       SenderService
}

func NewAppController(
	exchangeRateService RateService,
	emailSubscriptionService SubscriptionService,
	emailSenderService SenderService,
) *AppController {
	return &AppController{
		ExchangeRateService:      exchangeRateService,
		EmailSubscriptionService: emailSubscriptionService,
		EmailSenderService:       emailSenderService,
	}
}

func (ac *AppController) GetRate(w http.ResponseWriter, r *http.Request) {
	exchangeRate, err := ac.ExchangeRateService.ExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if json.NewEncoder(w).Encode(exchangeRate) != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *AppController) SubscribeEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	err := ac.EmailSubscriptionService.Subscribe(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (ac *AppController) SendEmails(w http.ResponseWriter, r *http.Request) {
	exchangeRate, err := ac.ExchangeRateService.ExchangeRate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subscribers, err := ac.EmailSubscriptionService.Subscriptions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ac.EmailSenderService.SendExchangeRate(exchangeRate, subscribers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
