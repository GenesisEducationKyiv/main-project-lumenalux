package coingecko

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"gses2-app/pkg/config"
	"gses2-app/pkg/types"
)

// Represents data type for JSON response
type response map[string]map[string]float64

var (
	ErrHTTPRequestFailure       = errors.New("http request failure")
	ErrUnexpectedStatusSode     = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat = errors.New("unexpected response format")
)

const (
	_providerName = "CoingeckoRateProvider"
	_currencyFrom = "bitcoin"
	_currencyTo   = "uah"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type CoingeckoProvider struct {
	config     config.CoingeckoAPIConfig
	httpClient HTTPClient
	logFunc    func(providerName string, resp *http.Response)
}

func NewProvider(
	config config.CoingeckoAPIConfig,
	httpClient HTTPClient,
	logFunc func(providerName string, resp *http.Response),
) *CoingeckoProvider {
	return &CoingeckoProvider{
		config:     config,
		httpClient: httpClient,
		logFunc:    logFunc,
	}
}

func (p *CoingeckoProvider) Name() string {
	return _providerName
}

func (p *CoingeckoProvider) ExchangeRate() (types.Rate, error) {
	resp, err := p.requestAPI()
	if err != nil {
		return 0, err
	}

	return p.extractRateFromResponse(resp)
}

func (p *CoingeckoProvider) requestAPI() (*http.Response, error) {

	resp, err := p.httpClient.Get(p.config.URL)
	if err != nil {
		return nil, ErrHTTPRequestFailure
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%v: %d", ErrUnexpectedStatusSode, resp.StatusCode)
	}

	p.logFunc(_providerName, resp)
	return resp, nil
}

func (p *CoingeckoProvider) extractRateFromResponse(resp *http.Response) (types.Rate, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return p.config.DefaultRate, err
	}

	var data response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return p.config.DefaultRate, errors.Join(err, ErrUnexpectedResponseFormat)
	}

	exchangeRate, ok := data[_currencyFrom][_currencyTo]
	if !ok {
		return p.config.DefaultRate, ErrUnexpectedResponseFormat
	}

	return types.Rate(exchangeRate), nil
}
