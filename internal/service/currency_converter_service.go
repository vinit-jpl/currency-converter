package service

import (
	"currency-converter/internal/dto"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ThirdPartyResponse struct {
	Result   string             `json:"result"`
	BaseCode string             `json:"base_code"`
	Rates    map[string]float64 `json:"rates"`
}

func buildExchangeRateURL(exchangeRateURL, baseCurrency string) string {
	return fmt.Sprintf("%s/%s", exchangeRateURL, baseCurrency)
}
func (s *ConverterService) ConvertCurrency(req dto.ConvertCurrencyRequest) (*dto.ConvertCurrencyResponse, error) {

	// get the url
	url := s.Cfg.ThirdPartyURL
	fmt.Println(url)

	// set from and to currency to Upper case
	fromCurrency := strings.ToUpper(req.FromCurrency)
	toCurrency := strings.ToUpper(req.ToCurrency)

	finalURL := buildExchangeRateURL(url, fromCurrency)

	// configure a client with a timeout to prevent hanging connections
	client := &http.Client{Timeout: 10 * time.Second}

	// make the get request
	resp, err := client.Get(finalURL)
	if err != nil {
		return nil, errors.New("Failed to fetch the data")
	}

	defer resp.Body.Close()

	// check for successful http status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// decode the response
	var exchangeResp ThirdPartyResponse

	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return nil, fmt.Errorf("failed to decode the response %w", err)
	}

	rate, ok := exchangeResp.Rates[toCurrency]
	if !ok {
		return nil, fmt.Errorf("unsupported currency: %s", toCurrency)
	}

	return &dto.ConvertCurrencyResponse{
		ExchangeRate:         rate,
		TotalConvertedAmount: rate * req.Amount,
		ToCurrency:           toCurrency,
	}, nil
}

func (s *ConverterService) GetCurrencies() (*dto.CurrenciesResponse, error) {
	url := buildExchangeRateURL(s.Cfg.ThirdPartyURL, "USD")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var exchangeResp ThirdPartyResponse
	if err := json.NewDecoder(resp.Body).Decode(&exchangeResp); err != nil {
		return nil, err
	}

	codes := []string{"USD"}
	for code := range exchangeResp.Rates {
		codes = append(codes, strings.ToUpper(code))
	}

	sort.Strings(codes)

	options := make([]dto.CurrencyOption, 0, len(codes))
	for _, code := range codes {
		name := code
		switch code {
		case "USD":
			name = "US Dollar"
		case "EUR":
			name = "Euro"
		case "INR":
			name = "Indian Rupee"
		case "JPY":
			name = "Japanese Yen"
		case "GBP":
			name = "British Pound"
		}

		options = append(options, dto.CurrencyOption{
			Code: code,
			Name: name,
		})
	}

	return &dto.CurrenciesResponse{Currencies: options}, nil
}
