package kuna

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gses2-app/internal/core/port"
	"gses2-app/internal/repository/rate/rest"
)

var (
	ErrHTTPRequestFailure           = errors.New("http request failure")
	ErrUnexpectedStatusCode         = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat     = errors.New("unexpected response format")
	ErrUnexpectedExchangeRateFormat = errors.New("unexpected exchange rate format")
)

const (
	_providerName     = "KunaRateProvider"
	_firstItemIndex   = 0
	_minResponseItems = 9
	_rateIndex        = 7
)

type KunaAPIConfig struct {
	URL string `default:"https://api.kuna.io/v3/tickers?symbols=btcuah"`
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type KunaProvider struct {
	config KunaAPIConfig
}

func NewProvider(
	logger port.Logger,
	config KunaAPIConfig,
	httpClient HTTPClient,
) *rest.AbstractProvider {
	return rest.NewProvider(
		logger,
		&KunaProvider{
			config: config,
		},
		httpClient,
	)
}

func (p *KunaProvider) URL() string {
	return p.config.URL
}

func (p *KunaProvider) Name() string {
	return _providerName
}

func (p *KunaProvider) ExtractRate(resp *http.Response) (port.Rate, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 || len(data[_firstItemIndex]) < _minResponseItems {
		return 0, ErrUnexpectedResponseFormat
	}

	exchangeRate, ok := data[_firstItemIndex][_rateIndex].(float64)
	if !ok {
		return 0, ErrUnexpectedExchangeRateFormat
	}

	return port.Rate(exchangeRate), nil
}
