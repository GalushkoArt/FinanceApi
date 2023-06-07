package apiclient

import (
	"context"
	"errors"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var UnknownTwelveDataError = errors.New("unknown twelve data error")

type TwelveDataApiKeyError struct {
	Err error
}

func (e *TwelveDataApiKeyError) Is(target error) bool {
	_, ok := target.(*TwelveDataApiKeyError)
	if !ok {
		return false
	}
	return e.Error() == target.Error()
}

func (e *TwelveDataApiKeyError) Error() string {
	return "problem with api key: " + e.Err.Error()
}

type twelveDataClient struct {
	c      *http.Client
	apiKey string
	host   string
}

type TwelveDataClient interface {
	GetHistoricDataForSymbol(ctx context.Context, symbol string) (*TimeSeries, error)
}

func NewTwelveDataClient(apiKey string, apiHost string, clientTimout time.Duration) TwelveDataClient {
	matched, err := regexp.MatchString("http(s)?://[\\w\\-.]+(:\\d{1,5})?", apiHost)
	if !matched {
		log.Panic().Err(err).Msg("Invalid apiHost!")
	}
	return &twelveDataClient{
		c: &http.Client{
			Timeout:   clientTimout,
			Transport: NewRequestLoggingTransport(http.DefaultTransport),
		},
		apiKey: apiKey,
		host:   apiHost,
	}
}

func (c *twelveDataClient) getUrl(resource string, params *url.Values) string {
	u, _ := url.ParseRequestURI(c.host)
	u.Path = resource
	u.RawQuery = params.Encode()
	result := u.String()
	return result
}

func (c *twelveDataClient) GetHistoricDataForSymbol(ctx context.Context, symbol string) (*TimeSeries, error) {
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
		return nil, UnknownTwelveDataError
	}
	var results TimeSeries
	if err = json.Unmarshal(responseBody, &results); err != nil {
		return nil, err
	}
	if results.Status == "error" {
		if results.Code == 400 || results.Code == 404 {
			return &results, model.SymbolNotFound
		} else if results.Code == 401 || results.Code == 429 {
			return &results, &TwelveDataApiKeyError{Err: errors.New("message: " + results.Message)}
		}
		return &results, UnknownTwelveDataError
	}
	return &results, nil
}
