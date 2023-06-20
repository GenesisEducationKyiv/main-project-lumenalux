package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const EXCHANGE_API_REQUEST_URL = "https://api.kuna.io/v3/tickers?symbols=btcuah"

type ExchangeRateService interface {
	GetExchangeRate() (float32, error)
}

type HTTPClient interface {
	Get(url string) (*http.Response, error)
}

type ExchangeRateServiceImpl struct {
	httpClient HTTPClient
}

func NewExchangeRateService(client HTTPClient) ExchangeRateService {
	return &ExchangeRateServiceImpl{
		httpClient: client,
	}
}

func (s *ExchangeRateServiceImpl) GetExchangeRate() (float32, error) {
	resp, err := s.makeRequest()
	if err != nil {
		return 0, err
	}

	data, err := s.parseResponse(resp)
	if err != nil {
		return 0, err
	}

	return s.extractExchangeRate(data)
}

func (s *ExchangeRateServiceImpl) makeRequest() (*http.Response, error) {
	resp, err := s.httpClient.Get(EXCHANGE_API_REQUEST_URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (s *ExchangeRateServiceImpl) parseResponse(resp *http.Response) ([][]interface{}, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data [][]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (s *ExchangeRateServiceImpl) extractExchangeRate(data [][]interface{}) (float32, error) {

	if len(data) == 0 || len(data[0]) < 9 {
		return 0, fmt.Errorf("unexpected response format")
	}

	exchangeRate, ok := data[0][7].(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected exchange rate format")
	}

	return float32(exchangeRate), nil
}
