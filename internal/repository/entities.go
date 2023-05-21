package repository

import (
	"time"
)

type price struct {
	SymbolID int64     `db:"symbol_id"`
	Date     time.Time `db:"date"`
	Open     string    `db:"open"`
	High     string    `db:"high"`
	Low      string    `db:"low"`
	Close    string    `db:"close"`
	Volume   string    `db:"volume"`
}

type symbol struct {
	ID            int64  `db:"id"`
	Symbol        string `db:"symbol"`
	Name          string `db:"name"`
	SymbolType    string `db:"type"`
	Currency      string `db:"currency"`
	CurrencyBase  string `db:"currency_base"`
	CurrencyQuote string `db:"currency_quote"`
}

type exchange struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Country  string `db:"country"`
	Code     string `db:"code"`
	Timezone string `db:"timezone"`
}

type userEntity struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type refreshToken struct {
	ID        int64     `db:"id"`
	UserId    string    `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
}
