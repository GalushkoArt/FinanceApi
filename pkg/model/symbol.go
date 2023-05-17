package model

type Symbol struct {
	ID            int64      `json:"-"`
	Symbol        string     `json:"symbol,omitempty" binding:"required"`
	Name          string     `json:"name,omitempty"`
	Type          string     `json:"type,omitempty"`
	Currency      string     `json:"currency,omitempty"`
	CurrencyBase  string     `json:"currency_base,omitempty"`
	CurrencyQuote string     `json:"currency_quote,omitempty"`
	Exchanges     []Exchange `json:"exchanges,omitempty"`
	Values        []Price    `json:"values,omitempty"`
}

type Exchange struct {
	ID       int64  `json:"-"`
	Name     string `json:"name,omitempty"`
	Country  string `json:"country,omitempty"`
	Timezone string `json:"timezone,omitempty"`
	MicCode  string `json:"mic_code,omitempty"`
}

type Price struct {
	Date   string `json:"date,omitempty"`
	Open   string `json:"open,omitempty"`
	High   string `json:"high,omitempty"`
	Low    string `json:"low,omitempty"`
	Close  string `json:"close,omitempty"`
	Volume string `json:"volume,omitempty"`
}
