package repository

import (
	"FinanceApi/pkg/log"
	"FinanceApi/pkg/model"
	"FinanceApi/pkg/utils"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type symbolRepositoryPostgres struct {
	db *sqlx.DB
}

type SymbolRepository interface {
	Add(symbol model.Symbol) error
	GetBySymbol(symbolName string) (model.Symbol, error)
	GetAll() ([]model.Symbol, error)
	Update(symbol model.Symbol) error
	Delete(symbolName string) error
}

var SymbolNotFound = errors.New("symbol not found")

func NewSymbolRepository(db *sqlx.DB) SymbolRepository {
	return &symbolRepositoryPostgres{db: db}
}

const (
	symbolQuery            = `SELECT id, symbol, name, type, currency, currency_base, currency_quote FROM SYMBOL WHERE SYMBOL = $1`
	exchangeQuery          = `SELECT id, name, country, code, timezone FROM EXCHANGE WHERE NAME = $1`
	symbolsWithLatestPrice = `SELECT id, symbol, name, type, currency, currency_base, currency_quote, date, open, close, high, low, volume FROM v_latest_symbol_info`
	symbolWithLatestPrice  = `SELECT id, symbol, name, type, currency, currency_base, currency_quote, date, open, close, high, low, volume FROM v_latest_symbol_info where symbol = $1`
)

func (r *symbolRepositoryPostgres) Add(newSymbol model.Symbol) error {
	var stored symbol
	tx, err := r.db.Beginx()
	err = tx.Get(&stored, symbolQuery, newSymbol.Symbol)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Panic("Found error in symbol select!", err)
		}
		log.InfoF("New symbol %s reference not found!", newSymbol.Symbol)
		const symbolInsert = `INSERT INTO SYMBOL (symbol, name, type, currency, currency_base, currency_quote) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		row := tx.QueryRow(symbolInsert, newSymbol.Symbol, newSymbol.Name, newSymbol.Type, newSymbol.Currency, newSymbol.CurrencyBase, newSymbol.CurrencyQuote)
		if err := row.Scan(&stored.Id); err != nil {
			utils.PanicOnError(tx.Rollback())
			return err
		}
	}
	if len(newSymbol.Exchanges) > 0 {
		var storedExchange exchange
		for _, exchange := range newSymbol.Exchanges {
			err = tx.Get(&storedExchange, exchangeQuery, exchange.Name)
			if err != nil {
				if err != sql.ErrNoRows {
					log.Panic("Found error in exchange select!", err)
				}
				log.InfoF("New symbol exchange %s reference not found!", exchange.Name)
				const exchangeInsert = `INSERT INTO EXCHANGE (name, code, country, timezone) VALUES ($1, $2, $3, $4)`
				_, err := tx.Exec(exchangeInsert, exchange.Name, exchange.MicCode, exchange.Country, exchange.Timezone)
				if err != nil {
					utils.PanicOnError(tx.Rollback())
					return err
				}
			}
		}
	}
	for _, price := range newSymbol.Values {
		const priceInsert = `INSERT INTO PRICE(SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := tx.Exec(priceInsert, stored.Id, price.Date, price.Open, price.Close, price.High, price.Low, price.Volume)
		if err != nil {
			log.Warn("Fail on insert new price!", err)
			utils.PanicOnError(tx.Rollback())
			return err
		}
	}
	return tx.Commit()
}

func (r *symbolRepositoryPostgres) Update(newSymbol model.Symbol) error {
	var stored symbol
	tx, err := r.db.Beginx()
	err = tx.Get(&stored, symbolQuery, newSymbol.Symbol)
	if err != nil {
		utils.PanicOnError(tx.Rollback())
		log.InfoF("Cannot update %s symbol!\n%v", newSymbol.Symbol, err)
		if err == sql.ErrNoRows {
			return SymbolNotFound
		}
		return err
	}
	const symbolUpdate = `update symbol set symbol = $1, name = $2, type = $3, currency = $4, currency_base = $5, currency_quote = $6 where id = $7`
	_, err = tx.Exec(symbolUpdate, newSymbol.Symbol, newSymbol.Name, newSymbol.Type, newSymbol.Currency, newSymbol.CurrencyBase, newSymbol.CurrencyQuote, stored.Id)
	if err != nil {
		log.InfoF("Cannot update %s symbol!\n%v", newSymbol.Symbol, err)
		utils.PanicOnError(tx.Rollback())
		return err
	}
	if len(newSymbol.Exchanges) > 0 {
		var storedExchange exchange
		for _, exchange := range newSymbol.Exchanges {
			err = tx.Get(&storedExchange, exchangeQuery, exchange.Name)
			if err != nil {
				log.InfoF("Cannot update %s exchange!\n%v", exchange.Name, err)
				utils.PanicOnError(tx.Rollback())
				return err
			}
			const exchangeUpdate = `UPDATE EXCHANGE name = $1, code = $2, country = $3, timezone = $4 where id = $5`
			_, err := tx.Exec(exchangeUpdate, exchange.Name, exchange.MicCode, exchange.Country, exchange.Timezone, storedExchange.Id)
			if err != nil {
				log.InfoF("Fail on update %s exchange!\n%v", exchange.Name, err)
			}
		}
	}
	var storedPrice price
	for _, price := range newSymbol.Values {
		const priceQuery = `SELECT SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME FROM PRICE WHERE SYMBOL_ID = $1 AND DATE = $2`
		err = tx.Get(&storedPrice, priceQuery, stored.Id, price.Date)
		if err != nil {
			const priceInsert = `INSERT INTO PRICE(SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err := tx.Exec(priceInsert, stored.Id, price.Date, price.Open, price.Close, price.High, price.Low, price.Volume)
			if err != nil {
				log.Info("Fail on insert price!", err)
			}
		} else {
			const priceUpdate = `UPDATE PRICE SET OPEN = $1, CLOSE = $2, HIGH = $3, LOW = $4, VOLUME = $5 WHERE SYMBOL_ID = $6 AND DATE = $7`
			_, err := tx.Exec(priceUpdate, price.Open, price.Close, price.High, price.Low, price.Volume, stored.Id, price.Date)
			if err != nil {
				log.Info("Fail on insert price!", err)
			}
		}
	}
	return tx.Commit()
}

func (r *symbolRepositoryPostgres) Delete(symbolName string) error {
	const symbolDelete = `DELETE FROM symbol where SYMBOL = $1`
	result, err := r.db.Exec(symbolDelete, symbolName)
	if err != nil {
		log.InfoF("Cannot delete %s symbol!\n%v", symbolName, err)
		return err
	}
	affected, _ := result.RowsAffected()
	if affected < 1 {
		return SymbolNotFound
	}
	return err
}

func (r *symbolRepositoryPostgres) GetBySymbol(symbolName string) (model.Symbol, error) {
	rows, err := r.db.Query(symbolWithLatestPrice, symbolName)
	if err != nil {
		return model.Symbol{}, err
	}
	results, err := r.retrieveLatest(rows)
	if len(results) < 1 || err == sql.ErrNoRows {
		return model.Symbol{}, SymbolNotFound
	}
	if err != nil {
		return model.Symbol{}, err
	}
	result := results[0]
	var storedExchanges []exchange
	err = r.db.Select(&storedExchanges, exchangeQuery, result.ID)
	if err != nil {
		return result, err
	}
	result.Exchanges = make([]model.Exchange, 0, len(storedExchanges))
	for _, storedExchange := range storedExchanges {
		result.Exchanges = append(result.Exchanges, model.Exchange{
			Name:     storedExchange.Name,
			Country:  storedExchange.Country,
			MicCode:  storedExchange.Code,
			Timezone: storedExchange.Timezone})
	}
	return result, nil
}

func (r *symbolRepositoryPostgres) GetAll() ([]model.Symbol, error) {
	rows, err := r.db.Query(symbolsWithLatestPrice)
	if err != nil {
		return []model.Symbol{}, err
	}
	return r.retrieveLatest(rows)
}

func (r *symbolRepositoryPostgres) retrieveLatest(rows *sql.Rows) ([]model.Symbol, error) {
	result := make([]model.Symbol, 0)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error("Couldn't close row!", err)
		}
	}(rows)
	for rows.Next() {
		var s model.Symbol
		s.Values = make([]model.Price, 1)
		err := rows.Scan(&s.ID, &s.Symbol, &s.Name, &s.Type, &s.Currency, &s.CurrencyBase, &s.CurrencyQuote, &s.Values[0].Date, &s.Values[0].Open, &s.Values[0].Close, &s.Values[0].High, &s.Values[0].Low, &s.Values[0].Volume)
		if err != nil {
			log.Error("Error on scanning row!", err)
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
