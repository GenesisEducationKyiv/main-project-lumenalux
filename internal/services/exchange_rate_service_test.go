package services

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func TestGetExchangeRate(t *testing.T) {
	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			jsonResponse := `[["btcuah",106,2.60,107,5.07,-14,-0.14,106,0.917,107,105]]`
			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))
			return &http.Response{
				StatusCode: 200,
				Body:       r,
			}, nil
		},
	}

	service := NewExchangeRateService(mockClient)
	rate, err := service.GetExchangeRate()
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	expectedRate := float32(106.0)
	if rate != expectedRate {
		t.Errorf("Expected rate %f, got: %f", expectedRate, rate)
	}
}

func TestGetExchangeRate_Failed(t *testing.T) {
	mockClient := &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			return nil, errors.New("error")
		},
	}

	service := &ExchangeRateServiceImpl{httpClient: mockClient}

	_, err := service.GetExchangeRate()
	if err == nil || !strings.Contains(err.Error(), "error") {
		t.Errorf("Expected error, got: %v", err)
	}

	mockClient = &MockHTTPClient{
		GetFunc: func(url string) (*http.Response, error) {
			var jsonResponse string = ""
			r := io.NopCloser(bytes.NewReader([]byte(jsonResponse)))
			return &http.Response{
				StatusCode: 400,
				Body:       r,
			}, nil
		},
	}

	service = &ExchangeRateServiceImpl{httpClient: mockClient}
	_, err = service.GetExchangeRate()
	if err == nil {
		t.Errorf("Expected error, got: %v", err)
	}
}
