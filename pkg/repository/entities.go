package repository

import "time"

type price struct {
	SymbolId int64     `db:"symbol_id"`
	Date     time.Time `db:"date"`
	Open     string    `db:"open"`
	High     string    `db:"high"`
	Low      string    `db:"low"`
	Close    string    `db:"close"`
	Volume   string    `db:"volume"`
}

type symbol struct {
	Id            int64  `db:"id"`
	Symbol        string `db:"symbol"`
	Name          string `db:"name"`
	SymbolType    string `db:"type"`
	Currency      string `db:"currency"`
	CurrencyBase  string `db:"currency_base"`
	CurrencyQuote string `db:"currency_quote"`
}

type exchange struct {
	Id       int64  `db:"id"`
	Name     string `db:"name"`
	Country  string `db:"country"`
	Code     string `db:"code"`
	Timezone string `db:"timezone"`
}

type SymbolInfo struct {
	Id            int64     `db:"id"`
	Symbol        string    `db:"symbol"`
	Name          string    `db:"name"`
	SymbolType    string    `db:"type"`
	Currency      string    `db:"currency"`
	CurrencyBase  string    `db:"currency_base"`
	CurrencyQuote string    `db:"currency_quote"`
	Date          time.Time `db:"date"`
	Open          string    `db:"open"`
	High          string    `db:"high"`
	Low           string    `db:"low"`
	Close         string    `db:"close"`
	Volume        string    `db:"volume"`
}
