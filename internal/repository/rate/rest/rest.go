package rest

import (
	"errors"
	"fmt"
	"gses2-app/internal/core/port"
	"net/http"
)

var (
	ErrHTTPRequestFailure   = errors.New("http request failure")
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type Provider interface {
	URL() string
	Name() string
	ExtractRate(resp *http.Response) (port.Rate, error)
}

type AbstractProvider struct {
	logger         port.Logger
	actualProvider Provider
	httpClient     HTTPClient
}

func NewProvider(
	logger port.Logger,
	actualProvider Provider,
	httpClient HTTPClient,
) *AbstractProvider {
	return &AbstractProvider{
		logger:         logger,
		actualProvider: actualProvider,
		httpClient:     httpClient,
	}
}

func (ap *AbstractProvider) Name() string {
	return ap.actualProvider.Name()
}

func (ap *AbstractProvider) ExchangeRate() (port.Rate, error) {
	resp, err := ap.requestAPI()
	if err != nil {
		return 0, err
	}

	return ap.extractRateFromResponse(resp)
}

func (ap *AbstractProvider) requestAPI() (*http.Response, error) {
	resp, err := ap.httpClient.Get(ap.actualProvider.URL())
	if err != nil {
		return nil, ErrHTTPRequestFailure
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedStatusCode, resp.StatusCode)
	}

	ap.logger.Infof("%v - Response: %v", ap.actualProvider.Name(), resp)
	return resp, nil
}

func (ap *AbstractProvider) extractRateFromResponse(resp *http.Response) (port.Rate, error) {
	defer resp.Body.Close()
	return ap.actualProvider.ExtractRate(resp)
}
