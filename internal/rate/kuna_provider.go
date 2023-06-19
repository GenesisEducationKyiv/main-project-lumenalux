package rate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gses2-app/pkg/config"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type KunaProvider struct {
	httpClient HTTPClient
}

func NewKunaProvider(httpClient HTTPClient) *KunaProvider {
	return &KunaProvider{
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

	resp, err := p.httpClient.Get(config.Current().KunaAPI.Url)
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
		return 0, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, err
	}

	if len(data) == 0 || len(data[0]) < 9 {
		return 0, fmt.Errorf("unexpected response format")
	}

	exchangeRate, ok := data[0][7].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected exchange rate format")
	}

	return float32(exchangeRate), nil
}
