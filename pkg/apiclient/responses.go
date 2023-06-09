package apiclient

type TimeSeries struct {
	Meta    Meta          `json:"meta,omitempty"`
	Values  []SeriesValue `json:"values,omitempty"`
	Code    int           `json:"code,omitempty"`
	Message string        `json:"message,omitempty"`
	Status  string        `json:"status,omitempty"`
}

type Meta struct {
	Symbol           string `json:"symbol,omitempty"`
	Interval         string `json:"interval,omitempty"`
	Currency         string `json:"currency,omitempty"`
	CurrencyBase     string `json:"currency_base,omitempty"`
	CurrencyQuote    string `json:"currency_quote,omitempty"`
	ExchangeTimezone string `json:"exchange_timezone,omitempty"`
	Exchange         string `json:"exchange,omitempty"`
	MicCode          string `json:"mic_code,omitempty"`
	Type             string `json:"type,omitempty"`
}

type SeriesValue struct {
	Datetime string `json:"datetime,omitempty"`
	Open     string `json:"open,omitempty"`
	High     string `json:"high,omitempty"`
	Low      string `json:"low,omitempty"`
	Close    string `json:"close,omitempty"`
	Volume   string `json:"volume,omitempty"`
}
