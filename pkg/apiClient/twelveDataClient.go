package apiClient

import (
	"FinanceApi/pkg/config"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TwelveDataClient struct {
	c      *http.Client
	apiKey string
	host   string
}

func NewTwelveDataClient() *TwelveDataClient {
	return &TwelveDataClient{
		c: &http.Client{
			Timeout:   config.Conf.API.TwelveData.Timeout,
			Transport: NewRequestLoggingTransport(http.DefaultTransport),
		},
		apiKey: config.Conf.API.TwelveData.ApiKey,
		host:   config.Conf.API.TwelveData.Host,
	}
}

func (c *TwelveDataClient) getUrl(resource string, params *url.Values) string {
	u, _ := url.ParseRequestURI(c.host)
	u.Path = resource
	u.RawQuery = params.Encode()
	result := u.String()
	return result
}

func (c *TwelveDataClient) GetHistoricDataForSymbol(symbol string) (*TimeSeries, error) {
	params := url.Values{}
	params.Add("apikey", c.apiKey)
	params.Add("symbol", symbol)
	params.Add("interval", "1day")
	response, err := c.c.Get(c.getUrl("time_series", &params))
	if err != nil {
		return nil, err
	}
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		var result ErrorResponse
		if err = json.Unmarshal(responseBody, &result); err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("error within request: %d code, message: %s", result.Code, result.Message))
	}
	var results TimeSeries
	if err = json.Unmarshal(responseBody, &results); err != nil {
		return nil, err
	}
	if results.Status == "error" {
		return nil, errors.New("unknown error in response: " + string(responseBody))
	}
	return &results, nil
}
