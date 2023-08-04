package router

import (
	"net/http"
	"time"
)

type HTTPConfig struct {
	Port    string        `default:"8080"`
	Timeout time.Duration `default:"10s"`
}

type Controller interface {
	GetRate(w http.ResponseWriter, r *http.Request)
	SubscribeEmail(w http.ResponseWriter, r *http.Request)
	SendEmails(w http.ResponseWriter, r *http.Request)
}

type httpRouter struct {
	controller Controller
}

func NewHTTPRouter(controller Controller) *httpRouter {
	return &httpRouter{controller: controller}
}

func (router *httpRouter) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/rate", router.controller.GetRate)
	mux.HandleFunc("/api/subscribe", router.controller.SubscribeEmail)
	mux.HandleFunc("/api/sendEmails", router.controller.SendEmails)
}
