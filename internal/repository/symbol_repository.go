package repository

import (
	"context"
	"database/sql"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type symbolRepositoryPostgres struct {
	db *sqlx.DB
}

type SymbolRepository interface {
	Add(ctx context.Context, symbol model.Symbol) error
	GetBySymbol(ctx context.Context, name string) (model.Symbol, error)
	GetAll(ctx context.Context) ([]model.Symbol, error)
	Update(ctx context.Context, symbol model.UpdateSymbol) error
	Delete(ctx context.Context, symbolName string) error
}

func srLog(c context.Context, e *zerolog.Event) *zerolog.Event {
	return utils.LogRequest(c, e).Str("from", "symbolRepositoryPostgres")
}

func NewSymbolRepository(db *sqlx.DB) SymbolRepository {
	return &symbolRepositoryPostgres{db: db}
}

const (
	symbolQuery              = `SELECT id, symbol, name, type, currency, currency_base, currency_quote FROM SYMBOL WHERE SYMBOL = $1`
	exchangeQuery            = `SELECT id, name, country, code, timezone FROM EXCHANGE WHERE NAME = $1`
	exchangesBySymbolIdQuery = `SELECT E.id, E.name, E.country, E.code, E.timezone FROM EXCHANGE E JOIN symbol_exchange se ON se.EXCHANGE_ID = E.ID  WHERE SYMBOL_ID = $1`
	symbolExchangeInsert     = `INSERT INTO symbol_exchange (symbol_id, exchange_id) VALUES ($1, $2)`
	symbolsWithLatestPrice   = `SELECT id, symbol, name, type, currency, currency_base, currency_quote, date, open, close, high, low, volume FROM v_latest_symbol_info`
	symbolWithLatestPrice    = `SELECT id, symbol, name, type, currency, currency_base, currency_quote, date, open, close, high, low, volume FROM v_latest_symbol_info where symbol = $1`
)

func (r *symbolRepositoryPostgres) Add(ctx context.Context, newSymbol model.Symbol) error {
	var stored symbol
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		srLog(ctx, log.Error()).Err(err).Msg("Failed to begin transaction")
		return err
	}
	err = tx.Get(&stored, symbolQuery, newSymbol.Symbol)
	srLog(ctx, log.Debug()).Msgf("Searching for %s symbol!", newSymbol.Symbol)
	if err != nil {
		if err != sql.ErrNoRows {
			srLog(ctx, log.Error()).Stack().Err(err).Msg("Found error in symbol select!")
			return err
		}
		srLog(ctx, log.Info()).Msgf("New symbol %s reference not found! Trying to insert received values", newSymbol.Symbol)
		const symbolInsert = `INSERT INTO SYMBOL (symbol, name, type, currency, currency_base, currency_quote) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
		row := tx.QueryRow(symbolInsert, newSymbol.Symbol, newSymbol.Name, newSymbol.Type, newSymbol.Currency, newSymbol.CurrencyBase, newSymbol.CurrencyQuote)
		if err := row.Scan(&stored.ID); err != nil {
			utils.PanicOnError(tx.Rollback())
			return err
		}
	}
	if len(newSymbol.Exchanges) > 0 {
		var storedExchange exchange
		for _, exchange := range newSymbol.Exchanges {
			srLog(ctx, log.Debug()).Msgf("Searching for %s exchange!", exchange.Name)
			err = tx.Get(&storedExchange, exchangeQuery, exchange.Name)
			if err != nil {
				if err != sql.ErrNoRows {
					srLog(ctx, log.Error()).Stack().Err(err).Msg("Found error in exchange select!")
					return err
				}
				srLog(ctx, log.Info()).Msgf("New symbol exchange %s reference not found Trying to insert received values!", exchange.Name)
				const exchangeInsert = `INSERT INTO EXCHANGE (name, code, country, timezone) VALUES ($1, $2, $3, $4) RETURNING id`
				row := tx.QueryRow(exchangeInsert, exchange.Name, exchange.MicCode, exchange.Country, exchange.Timezone)
				if err := row.Scan(&storedExchange.ID); err != nil {
					srLog(ctx, log.Warn()).Err(err).Msg("Fail on insert new exchange!")
					utils.PanicOnError(tx.Rollback())
					return err
				}
			}
			srLog(ctx, log.Debug()).Msgf("Searching for %s symbol - %s exchange relation!", newSymbol.Symbol, exchange.Name)
			_, err := tx.Exec(symbolExchangeInsert, stored.ID, storedExchange.ID)
			if err != nil {
				srLog(ctx, log.Warn()).Err(err).Msg("Fail on insert symbol exchange relation!")
				utils.PanicOnError(tx.Rollback())
				return err
			}
		}
	}
	for _, price := range newSymbol.Values {
		srLog(ctx, log.Debug()).Msgf("Inserting price %+v for %s", price, newSymbol.Symbol)
		const priceInsert = `INSERT INTO PRICE(SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := tx.Exec(priceInsert, stored.ID, price.Date, price.Open, price.Close, price.High, price.Low, price.Volume)
		if err != nil {
			srLog(ctx, log.Warn()).Err(err).Msg("Fail on insert new price!")
			utils.PanicOnError(tx.Rollback())
			return err
		}
	}
	return tx.Commit()
}

func (r *symbolRepositoryPostgres) Update(ctx context.Context, newSymbol model.UpdateSymbol) error {
	var stored symbol
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		srLog(ctx, log.Error()).Err(err).Msg("Failed to begin transaction")
		return err
	}
	srLog(ctx, log.Debug()).Msgf("Searching for %s symbol!", newSymbol.Symbol)
	err = tx.Get(&stored, symbolQuery, newSymbol.Symbol)
	if err != nil {
		srLog(ctx, log.Info()).Err(err).Msgf("Cannot update %s symbol!", newSymbol.Symbol)
		utils.PanicOnError(tx.Rollback())
		if err == sql.ErrNoRows {
			return model.SymbolNotFound
		}
		return err
	}
	srLog(ctx, log.Debug()).Msgf("Updating %s symbol!", newSymbol.Symbol)
	updatedSymbol := updatedSymbol(stored, newSymbol)
	const symbolUpdate = `update symbol set symbol = $1, name = $2, type = $3, currency = $4, currency_base = $5, currency_quote = $6 where id = $7`
	_, err = tx.Exec(symbolUpdate, updatedSymbol.Symbol, updatedSymbol.Name, updatedSymbol.SymbolType, updatedSymbol.Currency, updatedSymbol.CurrencyBase, updatedSymbol.CurrencyQuote, stored.ID)
	if err != nil {
		srLog(ctx, log.Info()).Err(err).Msgf("Cannot update %s symbol!", newSymbol.Symbol)
		utils.PanicOnError(tx.Rollback())
		return err
	}
	if len(newSymbol.Exchanges) > 0 {
		var storedExchange exchange
		for _, exchange := range newSymbol.Exchanges {
			err = tx.Get(&storedExchange, exchangeQuery, exchange.Name)
			if err != nil {
				srLog(ctx, log.Info()).Err(err).Msgf("Cannot update %s exchange!", exchange.Name)
				utils.PanicOnError(tx.Rollback())
				return err
			}
			const exchangeUpdate = `UPDATE EXCHANGE SET name = $1, code = $2, country = $3, timezone = $4 where id = $5`
			_, err := tx.Exec(exchangeUpdate, exchange.Name, exchange.MicCode, exchange.Country, exchange.Timezone, storedExchange.ID)
			if err != nil {
				srLog(ctx, log.Info()).Err(err).Msgf("Fail on update %s exchange!", exchange.Name)
			}
		}
	}
	var storedPrice price
	for _, price := range newSymbol.Values {
		const priceQuery = `SELECT SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME FROM PRICE WHERE SYMBOL_ID = $1 AND DATE = $2`
		err = tx.Get(&storedPrice, priceQuery, stored.ID, price.Date)
		if err != nil {
			const priceInsert = `INSERT INTO PRICE(SYMBOL_ID, DATE, OPEN, CLOSE, HIGH, LOW, VOLUME) VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err := tx.Exec(priceInsert, stored.ID, price.Date, price.Open, price.Close, price.High, price.Low, price.Volume)
			if err != nil {
				srLog(ctx, log.Info()).Err(err).Msg("Fail on insert price!")
			}
		} else {
			const priceUpdate = `UPDATE PRICE SET OPEN = $1, CLOSE = $2, HIGH = $3, LOW = $4, VOLUME = $5 WHERE SYMBOL_ID = $6 AND DATE = $7`
			_, err := tx.Exec(priceUpdate, price.Open, price.Close, price.High, price.Low, price.Volume, stored.ID, price.Date)
			if err != nil {
				srLog(ctx, log.Info()).Err(err).Msg("Fail on insert price!")
			}
		}
	}
	return tx.Commit()
}

func updatedSymbol(origin symbol, new model.UpdateSymbol) symbol {
	if new.Name != nil {
		origin.Name = *new.Name
	}
	if new.Type != nil {
		origin.SymbolType = *new.Type
	}
	if new.Currency != nil {
		origin.Currency = *new.Currency
	}
	if new.CurrencyBase != nil {
		origin.CurrencyBase = *new.CurrencyBase
	}
	if new.CurrencyQuote != nil {
		origin.CurrencyQuote = *new.CurrencyQuote
	}
	return origin
}

func (r *symbolRepositoryPostgres) Delete(ctx context.Context, symbolName string) error {
	const symbolDelete = `DELETE FROM symbol where SYMBOL = $1`
	result, err := r.db.ExecContext(ctx, symbolDelete, symbolName)
	if err != nil {
		srLog(ctx, log.Info()).Err(err).Msgf("Cannot delete %s symbol!", symbolName)
		return err
	}
	affected, _ := result.RowsAffected()
	if affected < 1 {
		return model.SymbolNotFound
	}
	return err
}

func (r *symbolRepositoryPostgres) GetBySymbol(ctx context.Context, symbolName string) (model.Symbol, error) {
	rows, err := r.db.QueryContext(ctx, symbolWithLatestPrice, symbolName)
	if err != nil {
		return model.Symbol{}, err
	}
	results, err := r.retrieveLatest(ctx, rows)
	if len(results) < 1 || err == sql.ErrNoRows {
		return model.Symbol{}, model.SymbolNotFound
	}
	if err != nil {
		return model.Symbol{}, err
	}
	result := results[0]
	var storedExchanges []exchange
	err = r.db.SelectContext(ctx, &storedExchanges, exchangesBySymbolIdQuery, result.ID)
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

func (r *symbolRepositoryPostgres) GetAll(ctx context.Context) ([]model.Symbol, error) {
	rows, err := r.db.QueryContext(ctx, symbolsWithLatestPrice)
	if err != nil {
		return []model.Symbol{}, err
	}
	return r.retrieveLatest(ctx, rows)
}

func (r *symbolRepositoryPostgres) retrieveLatest(ctx context.Context, rows *sql.Rows) ([]model.Symbol, error) {
	result := make([]model.Symbol, 0)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			srLog(ctx, log.Error()).Err(err).Msg("Couldn't close row!")
		}
	}(rows)
	for rows.Next() {
		var s model.Symbol
		s.Values = make([]model.Price, 1)
		err := rows.Scan(&s.ID, &s.Symbol, &s.Name, &s.Type, &s.Currency, &s.CurrencyBase, &s.CurrencyQuote, &s.Values[0].Date, &s.Values[0].Open, &s.Values[0].Close, &s.Values[0].High, &s.Values[0].Low, &s.Values[0].Volume)
		if err != nil {
			srLog(ctx, log.Error()).Err(err).Msg("Error on scanning row!")
		}
		result = append(result, s)
	}
	return result, rows.Err()
}
