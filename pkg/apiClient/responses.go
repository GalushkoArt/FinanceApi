package apiClient

type TimeSeries struct {
	Meta struct {
		Symbol           string `json:"symbol"`
		Interval         string `json:"interval"`
		Currency         string `json:"currency"`
		CurrencyBase     string `json:"currency_base"`
		CurrencyQuote    string `json:"currency_quote"`
		ExchangeTimezone string `json:"exchange_timezone"`
		Exchange         string `json:"exchange"`
		MicCode          string `json:"mic_code"`
		Type             string `json:"type"`
	} `json:"meta"`
	Values []struct {
		Datetime string `json:"datetime"`
		Open     string `json:"open"`
		High     string `json:"high"`
		Low      string `json:"low"`
		Close    string `json:"close"`
		Volume   string `json:"volume"`
	} `json:"values"`
	Status string `json:"status"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
