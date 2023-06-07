package apiclient

import (
	"context"
	"errors"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetHistoricDataForSymbol(t *testing.T) {
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			client := twelveDataClient{host: "http://localhost", apiKey: "test",
				c: utils.MockClient(func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.receivedCode,
						Body:       utils.BodyFromStruct(tt.expected),
					}, tt.transportError
				})}
			timeSeries, returnError := client.GetHistoricDataForSymbol(context.TODO(), "TEST")
			assert.Equal(t, tt.expected, timeSeries, "TimeSeries should be equal")
			assert.ErrorIs(t, returnError, tt.expectedError, "Error should be equal")
		})
	}
}

var transportError = errors.New("transport error")

var testData = []struct {
	name           string
	expected       *TimeSeries
	receivedCode   int
	transportError error
	expectedError  error
}{
	{utils.TestName("Positive test"),
		&TimeSeries{
			Meta: Meta{
				Symbol:           "TEST",
				Interval:         "1day",
				Currency:         "USD",
				ExchangeTimezone: "America/New_York",
				Exchange:         "NASDAQ",
				MicCode:          "XNGS",
				Type:             "Common Stock",
			},
			Values: []SeriesValue{
				{
					Datetime: "2023-06-02",
					Open:     "181.00000",
					High:     "181.78000",
					Low:      "179.67380",
					Close:    "179.80000",
					Volume:   "16538045",
				},
				{
					Datetime: "2023-06-01",
					Open:     "177.70000",
					High:     "180.12000",
					Low:      "176.92999",
					Close:    "180.09000",
					Volume:   "68813900",
				},
			},
			Status: "ok",
		},
		200,
		nil,
		nil,
	},
	{utils.TestName("Error from TwelveData 400"), &TimeSeries{Code: 400, Status: "error", Message: "Some error"}, 200, nil, model.SymbolNotFound},
	{utils.TestName("Error from TwelveData 401"), &TimeSeries{Code: 401, Status: "error", Message: "API key is invalid"}, 200, nil, &TwelveDataApiKeyError{errors.New("message: API key is invalid")}},
	{utils.TestName("Error from TwelveData 404"), &TimeSeries{Code: 404, Status: "error", Message: "Not Found"}, 200, nil, model.SymbolNotFound},
	{utils.TestName("Error from TwelveData 429"), &TimeSeries{Code: 429, Status: "error", Message: "Overuse"}, 200, nil, &TwelveDataApiKeyError{errors.New("message: Overuse")}},
	{utils.TestName("Error from TwelveData 500"), &TimeSeries{Code: 500, Status: "error", Message: "Some error"}, 200, nil, UnknownTwelveDataError},
	{utils.TestName("Empty body"), &TimeSeries{}, 200, nil, nil},
	{utils.TestName("Status code 404"), nil, 404, nil, UnknownTwelveDataError},
	{utils.TestName("Transport Error"), nil, 200, transportError, transportError},
}
