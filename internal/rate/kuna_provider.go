package rate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gses2-app/pkg/config"
)

const (
	firstItemIndex   = 0
	minResponseItems = 9
	rateIndex        = 7
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type KunaProvider struct {
	config     config.KunaAPIConfig
	httpClient HTTPClient
}

func NewKunaProvider(config config.KunaAPIConfig, httpClient HTTPClient) *KunaProvider {
	return &KunaProvider{
		config:     config,
		httpClient: httpClient,
	}
}

func (p *KunaProvider) ExchangeRate() (float32, error) {
	resp, err := p.requestAPI()
	if err != nil {
		return 0, err
	}

	return p.extractRateFromResponse(resp)
}

func (p *KunaProvider) requestAPI() (*http.Response, error) {

	resp, err := p.httpClient.Get(p.config.URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (p *KunaProvider) extractRateFromResponse(resp *http.Response) (float32, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return p.config.DefaltRate, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return p.config.DefaltRate, err
	}

	if len(data) == 0 || len(data[firstItemIndex]) < minResponseItems {
		return p.config.DefaltRate, fmt.Errorf("unexpected response format")
	}

	exchangeRate, ok := data[firstItemIndex][rateIndex].(float64)
	if !ok {
		return p.config.DefaltRate, fmt.Errorf("unexpected exchange rate format")
	}

	return float32(exchangeRate), nil
}
