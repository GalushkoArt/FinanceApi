package apiClient

import (
	"FinanceApi/internal/model"
	"FinanceApi/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"time"
)

type TwelveDataClient struct {
	c      *http.Client
	apiKey string
	host   string
}

func NewTwelveDataClient(apiKey string, apiHost string, clientTimout time.Duration) *TwelveDataClient {
	return &TwelveDataClient{
		c: &http.Client{
			Timeout:   clientTimout,
			Transport: NewRequestLoggingTransport(http.DefaultTransport),
		},
		apiKey: apiKey,
		host:   apiHost,
	}
}

func (c *TwelveDataClient) getUrl(resource string, params *url.Values) string {
	u, _ := url.ParseRequestURI(c.host)
	u.Path = resource
	u.RawQuery = params.Encode()
	result := u.String()
	return result
}

func (c *TwelveDataClient) GetHistoricDataForSymbol(ctx context.Context, symbol string) (*TimeSeries, error) {
	params := url.Values{}
	params.Add("apikey", c.apiKey)
	params.Add("symbol", symbol)
	params.Add("interval", "1day")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getUrl("time_series", &params), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.c.Do(request)
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		utils.LogRequest(ctx, log.Error()).
			Str("from", "twelveDataClient").
			Str("response", string(responseBody)).
			Int("response_code", response.StatusCode).
			Msg("Unknown status code from TwelveData API")
		return nil, errors.New("unknown integration fail")
	}
	var results TimeSeries
	if err = json.Unmarshal(responseBody, &results); err != nil {
		return nil, err
	}
	if results.Status == "error" {
		utils.LogRequest(ctx, log.Error()).
			Str("from", "twelveDataClient").
			Interface("response", results).
			Int("response_code", response.StatusCode).
			Msg("Error received from TwelveData API")
		if results.Code == 400 || results.Code == 404 {
			return nil, model.SymbolNotFound
		}
		return nil, errors.New("unknown error in response from api")
	}
	return &results, nil
}
