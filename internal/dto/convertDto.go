package dto

type (
	ConvertCurrencyRequest struct {
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
		Amount       float64 `json:"amount"`
	}

	ConvertCurrencyResponse struct {
		ExchangeRate         float64 `json:"exchange_rate"`
		TotalConvertedAmount float64 `json:"total_converted_amount"`
		ToCurrency           string  `json:"to_currency"`
	}

	CurrencyOption struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}

	CurrenciesResponse struct {
		Currencies []CurrencyOption `json:"currencies"`
	}
)
