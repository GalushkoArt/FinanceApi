package apiClient

import (
	"FinanceApi/pkg/log"
	"fmt"
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
	logString := fmt.Sprintf("%s request to %s", r.Method, r.URL.String())
	if len(r.Header) > 0 {
		logString += fmt.Sprintf(" Header: %+v", r.Header)
	}
	if r.Body != nil {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		if len(body) > 0 {
			logString += " Body: " + string(body)
		}
	}
	log.Info(logString)
	return t.next.RoundTrip(r)
}
