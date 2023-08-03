package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"gses2-app/internal/rate"
	"gses2-app/internal/subscription"
	"gses2-app/internal/user/repository"
)

type SenderService interface {
	SendExchangeRate(rate rate.Rate, subscribers ...repository.User) error
}

type RateService interface {
	ExchangeRate() (rate rate.Rate, err error)
}

type SubscriptionService interface {
	Subscribe(subscriber *repository.User) error
	Subscriptions() (subscribers []repository.User, err error)
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

	if err = json.NewEncoder(w).Encode(exchangeRate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ac *AppController) SubscribeEmail(w http.ResponseWriter, r *http.Request) {
	subscriber := &repository.User{Email: r.FormValue("email")}
	err := ac.EmailSubscriptionService.Subscribe(subscriber)

	if errors.Is(err, subscription.ErrAlreadySubscribed) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = ac.EmailSenderService.SendExchangeRate(
		exchangeRate,
		subscribers...,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
