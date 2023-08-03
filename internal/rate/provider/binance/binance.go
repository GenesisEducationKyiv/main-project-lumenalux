package binance

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"gses2-app/internal/rate/provider"
	"gses2-app/pkg/config"
	"gses2-app/pkg/types"
)

var (
	ErrHTTPRequestFailure           = errors.New("http request failure")
	ErrUnexpectedStatusCode         = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat     = errors.New("unexpected response format")
	ErrUnexpectedExchangeRateFormat = errors.New("unexpected exchange rate format")
)

const (
	_providerName     = "BinanceRateProvider"
	_firstItemIndex   = 0
	_minResponseItems = 5
	_rateIndex        = 4
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type BinanceProvider struct {
	config config.BinanceAPIConfig
}

func NewProvider(
	config config.BinanceAPIConfig,
	httpClient HTTPClient,
	logFunc func(providerName string, resp *http.Response),
) *provider.AbstractProvider {
	return provider.NewProvider(
		&BinanceProvider{
			config: config,
		},
		httpClient,
		logFunc,
	)
}

func (p *BinanceProvider) URL() string {
	return p.config.URL
}

func (p *BinanceProvider) Name() string {
	return _providerName
}

func (p *BinanceProvider) ExtractRate(resp *http.Response) (types.Rate, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return p.config.DefaultRate, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return p.config.DefaultRate, err
	}

	if len(data) == 0 || len(data[_firstItemIndex]) < _minResponseItems {
		return p.config.DefaultRate, ErrUnexpectedResponseFormat
	}

	exchangeRate, ok := data[_firstItemIndex][_rateIndex].(string)
	if !ok {
		return p.config.DefaultRate, ErrUnexpectedExchangeRateFormat
	}

	rate, err := strconv.ParseFloat(exchangeRate, 64)
	if err != nil {
		return p.config.DefaultRate, errors.Join(err, ErrUnexpectedExchangeRateFormat)
	}

	return types.Rate(rate), nil
}
