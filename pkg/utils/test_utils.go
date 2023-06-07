package utils

import (
	"bytes"
	"fmt"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"
)

func MockClient(transport func(r *http.Request) (*http.Response, error)) *http.Client {
	return &http.Client{
		Transport: &mockTransport{
			mock: transport,
		},
	}
}

type mockTransport struct {
	http.RoundTripper
	mock func(r *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.mock(r)
}

func BodyFromStruct(s interface{}) io.ReadCloser {
	jsonString, err := json.Marshal(s)
	PanicOnError(err)
	return io.NopCloser(bytes.NewReader(jsonString))
}

func DeleteRequest(url string, body interface{}, noContentType bool, headers ...map[string]string) *http.Request {
	return Request(http.MethodDelete, url, body, noContentType, headers...)
}

func GetRequest(url string, headers ...map[string]string) *http.Request {
	return Request(http.MethodGet, url, nil, false, headers...)
}

func PostRequest(url string, body interface{}, noContentType bool, headers ...map[string]string) *http.Request {
	return Request(http.MethodPost, url, body, noContentType, headers...)
}

func PutRequest(url string, body interface{}, noContentType bool, headers ...map[string]string) *http.Request {
	return Request(http.MethodPut, url, body, noContentType, headers...)
}

func Request(method string, url string, body interface{}, noContentType bool, headers ...map[string]string) *http.Request {
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(body)
	req, err := http.NewRequest(method, url, &b)
	PanicOnError(err)
	if !noContentType && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	}
	if len(headers) != 0 {
		for key, value := range headers[0] {
			req.Header.Set(key, value)
		}
	}
	return req
}

var recorderPool = sync.Pool{
	New: func() interface{} {
		return httptest.NewRecorder()
	},
}

func SetCookie(req *http.Request, cookie *http.Cookie) {
	recorder := recorderPool.Get().(*httptest.ResponseRecorder)
	http.SetCookie(recorder, cookie)
	req.Header.Set("Cookie", recorder.Header().Get("Set-Cookie"))
	recorder.Header().Del("Set-Cookie")
	recorderPool.Put(recorder)
}

// TestName returns the name of the test with the file and line number
// Use it for test table names
func TestName(name string) string {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return name
	}
	return fmt.Sprintf("%s:%d/%s", file[strings.LastIndex(file, "/")+1:], line, name)
}

func CommonResponseAssertions(t *testing.T, response *http.Response, err error, expectedStatusCode int, expectedBody interface{}) {
	assert.NoError(t, err)
	assert.Equalf(t, expectedStatusCode, response.StatusCode, "expected status code %d, got %d", expectedStatusCode, response.StatusCode)
	expectedJson, err := json.Marshal(expectedBody)
	PanicOnError(err)
	responseJson, err := io.ReadAll(response.Body)
	PanicOnError(err)
	assert.JSONEq(t, string(expectedJson), string(responseJson))
}

func TestAuthMiddleware(c *fiber.Ctx) error {
	c.Locals("role", model.Role(c.GetReqHeaders()["Role"]))
	return c.Next()
}
