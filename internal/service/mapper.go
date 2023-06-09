package service

import (
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/pkg/apiclient"
)

func timeSeriesToModel(timeSeries apiclient.TimeSeries) model.Symbol {
	return model.Symbol{
		Symbol:        timeSeries.Meta.Symbol,
		Type:          timeSeries.Meta.Type,
		Currency:      timeSeries.Meta.Currency,
		CurrencyBase:  timeSeries.Meta.CurrencyBase,
		CurrencyQuote: timeSeries.Meta.CurrencyQuote,
		Exchanges: []model.Exchange{{
			Name:     timeSeries.Meta.Exchange,
			Timezone: timeSeries.Meta.ExchangeTimezone,
			MicCode:  timeSeries.Meta.MicCode,
		}},
		Values: []model.Price{{
			Date:   timeSeries.Values[0].Datetime,
			Open:   timeSeries.Values[0].Open,
			Close:  timeSeries.Values[0].Close,
			High:   timeSeries.Values[0].High,
			Low:    timeSeries.Values[0].Low,
			Volume: timeSeries.Values[0].Volume,
		}},
	}
}
