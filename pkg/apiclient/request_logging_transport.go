package apiclient

import (
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

type RequestLoggingTransport struct {
	next http.RoundTripper
}

func NewRequestLoggingTransport(next http.RoundTripper) *RequestLoggingTransport {
	return &RequestLoggingTransport{next: next}
}

func (t *RequestLoggingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	response, err := t.next.RoundTrip(r)
	responseBody, errRead := io.ReadAll(response.Body)
	if errRead != nil {
		responseBody = []byte("")
	}
	logEvent := log.Info()
	if response.StatusCode >= 400 {
		logEvent = log.Error()
	}
	utils.LogRequest(r.Context(), logEvent).
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Str("from", "apiClientRequestLogger").
		Interface("request_header", r.Header).
		Str("status", response.Status).
		Str("response_body", string(responseBody)).
		Msg("Request to " + r.Host)
	return response, err
}
