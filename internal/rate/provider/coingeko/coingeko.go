package coingeko

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gses2-app/internal/rate/provider"
	"gses2-app/pkg/config"
	"gses2-app/pkg/types"
)

// Represents data type for JSON response
type Response struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

var (
	ErrHTTPRequestFailure       = errors.New("http request failure")
	ErrUnexpectedStatusCode     = errors.New("unexpected status code")
	ErrUnexpectedResponseFormat = errors.New("unexpected response format")
)

const (
	_providerName = "CoingeckoRateProvider"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type CoingeckoProvider struct {
	config config.CoingeckoAPIConfig
}

func NewProvider(
	config config.CoingeckoAPIConfig,
	httpClient HTTPClient,
) *provider.AbstractProvider {
	return provider.NewProvider(
		&CoingeckoProvider{
			config: config,
		},
		httpClient,
	)
}

func (p *CoingeckoProvider) URL() string {
	return p.config.URL
}

func (p *CoingeckoProvider) Name() string {
	return _providerName
}

func (p *CoingeckoProvider) ExtractRate(resp *http.Response) (types.Rate, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, errors.Join(err, ErrUnexpectedResponseFormat)
	}

	return types.Rate(data.Bitcoin.UAH), nil
}
