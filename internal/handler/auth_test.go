package handler

import (
	"errors"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/service"
	"github.com/galushkoart/finance-api/mock"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

//go:generate echo $PWD - $GOFILE
//go:generate mockgen -package mock -destination ../../mock/auth_service_mock.go -source=../service/auth_service.go AuthService

func TestSignUp(t *testing.T) {
	mockService := mock.NewMockAuthService(gomock.NewController(t))
	app := setupFiberTest(&Handler{ah: authHandler{service: mockService}})
	for _, td := range signUpTestData {
		t.Run(td.name, func(t *testing.T) {
			if !td.wrongBody && !td.wrongContentType {
				mockService.EXPECT().SignUp(gomock.Any(), td.body).Return(td.serviceError)
			}
			response, err := app.Test(utils.PostRequest("/auth/signup", td.body, td.wrongContentType))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var signUpTestData = []struct {
	name             string
	body             model.SignUp
	expectedCode     int
	serviceError     error
	wrongContentType bool
	wrongBody        bool
	expectedResponse CommonResponse
}{
	{
		name:             utils.TestName("successful sign-up"),
		body:             model.SignUp{Username: "test", Email: "test@email.com", Password: "qwerty"},
		expectedCode:     200,
		expectedResponse: CommonResponse{Message: "success", Code: 200},
	},
	{
		name:             utils.TestName("wrong content type"),
		body:             model.SignUp{Username: "test", Email: "test@email.com", Password: "qwerty"},
		expectedCode:     400,
		wrongContentType: true,
		expectedResponse: CommonResponse{Message: "Wrong content type", Code: 400},
	},
	{
		name:             utils.TestName("wrong body"),
		body:             model.SignUp{Username: "te", Email: "not@an_email"},
		expectedCode:     400,
		wrongBody:        true,
		expectedResponse: CommonResponse{Message: "Wrong body", Code: 400, AuthErrors: []*model.AuthError{{Field: "Username", Rule: "min"}, {Field: "Email", Rule: "email"}, {Field: "Password", Rule: "min"}}},
	},
	{
		name:             utils.TestName("user already exists"),
		body:             model.SignUp{Username: "test", Email: "test@email.com", Password: "qwerty"},
		expectedCode:     400,
		serviceError:     service.UserAlreadyExists,
		expectedResponse: CommonResponse{Message: "User with such username or email already exists", Code: 400},
	},
	{
		name:             utils.TestName("wrong credentials"),
		body:             model.SignUp{Username: "test", Email: "test@email.com", Password: "qwerty"},
		expectedCode:     500,
		serviceError:     errors.New("failed to sign up"),
		expectedResponse: CommonResponse{Message: "Failed to register", Code: 500},
	},
}

func TestSignIn(t *testing.T) {
	mockService := mock.NewMockAuthService(gomock.NewController(t))
	app := setupFiberTest(&Handler{ah: authHandler{service: mockService}})
	for _, td := range signInTestData {
		t.Run(td.name, func(t *testing.T) {
			if !td.wrongBody && !td.wrongContentType {
				mockService.EXPECT().SignIn(gomock.Any(), td.body).Return(td.jwtToken, td.refreshToken, td.cookieExpiry, td.serviceError)
			}
			response, err := app.Test(utils.PostRequest("/auth/signin", td.body, td.wrongContentType))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
			if td.serviceError == nil && !td.wrongBody && !td.wrongContentType {
				cookie := *response.Cookies()[0]
				assert.Equal(t, cookie.HttpOnly, true, "cookie should be HttpOnly")
				assert.Equal(t, cookie.Expires.UTC().Unix(), td.cookieExpiry.UTC().Unix(), "cookie should expire as expected")
				assert.Equalf(t, cookie.Name, refreshTokenCookie, "cookie should have name '%s'", refreshTokenCookie)
				assert.Equal(t, cookie.Value, td.refreshToken, "cookie should have refresh token as value")
			}
		})
	}
}

var signInTestData = []struct {
	name             string
	body             model.SignIn
	expectedCode     int
	jwtToken         string
	refreshToken     string
	cookieExpiry     time.Time
	serviceError     error
	wrongContentType bool
	wrongBody        bool
	expectedResponse interface{}
}{
	{
		name:             utils.TestName("successful sign-in"),
		body:             model.SignIn{Login: "test", Password: "qwerty"},
		expectedCode:     200,
		jwtToken:         "jwt_token",
		refreshToken:     "refresh_token",
		cookieExpiry:     time.Now().Add(24 * time.Hour),
		expectedResponse: model.SuccessfulAuthentication{Token: "jwt_token"},
	},
	{
		name:             utils.TestName("wrong content type"),
		body:             model.SignIn{Login: "test", Password: "qwerty"},
		expectedCode:     400,
		wrongContentType: true,
		expectedResponse: CommonResponse{Message: "Wrong content type", Code: 400},
	},
	{
		name:             utils.TestName("wrong body"),
		body:             model.SignIn{Login: "ts"},
		expectedCode:     400,
		wrongBody:        true,
		expectedResponse: CommonResponse{Message: "Wrong body", Code: 400, AuthErrors: []*model.AuthError{{Field: "Login", Rule: "min"}, {Field: "Password", Rule: "min"}}},
	},
	{
		name:             utils.TestName("wrong credentials"),
		body:             model.SignIn{Login: "test", Password: "qwerty"},
		expectedCode:     401,
		serviceError:     model.UserNotFound,
		expectedResponse: CommonResponse{Message: "Wrong credentials", Code: 401},
	},
	{
		name:             utils.TestName("wrong credentials"),
		body:             model.SignIn{Login: "test", Password: "qwerty"},
		expectedCode:     401,
		serviceError:     model.UserNotFound,
		expectedResponse: CommonResponse{Message: "Wrong credentials", Code: 401},
	},
	{
		name:             utils.TestName("failed to sign in"),
		body:             model.SignIn{Login: "test", Password: "qwerty"},
		expectedCode:     500,
		serviceError:     errors.New("failed to sign in"),
		expectedResponse: CommonResponse{Message: "Failed to sign in", Code: 500},
	},
}

func TestRefresh(t *testing.T) {
	mockService := mock.NewMockAuthService(gomock.NewController(t))
	app := setupFiberTest(&Handler{ah: authHandler{service: mockService}})
	for _, td := range refreshTestData {
		t.Run(td.name, func(t *testing.T) {
			if !td.emptyRefreshToken {
				mockService.EXPECT().RefreshToken(gomock.Any(), td.refreshToken).Return(td.jwtToken, td.newRefreshToken, td.cookieExpiry, td.serviceError)
			}
			request := utils.GetRequest("/auth/refresh")
			utils.SetCookie(request, &http.Cookie{Name: refreshTokenCookie, Value: td.refreshToken})
			response, err := app.Test(request)
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
			if td.serviceError == nil && !td.emptyRefreshToken {
				cookie := *response.Cookies()[0]
				assert.Equal(t, cookie.HttpOnly, true, "cookie should be HttpOnly")
				assert.Equal(t, cookie.Expires.UTC().Unix(), td.cookieExpiry.UTC().Unix(), "cookie should expire as expected")
				assert.Equalf(t, cookie.Name, refreshTokenCookie, "cookie should have name '%s'", refreshTokenCookie)
				assert.Equalf(t, cookie.Value, td.newRefreshToken, "cookie value should match")
			}
		})
	}
}

var refreshTestData = []struct {
	name              string
	refreshToken      string
	expectedCode      int
	jwtToken          string
	newRefreshToken   string
	cookieExpiry      time.Time
	serviceError      error
	emptyRefreshToken bool
	expectedResponse  interface{}
}{
	{
		name:             utils.TestName("successful refresh"),
		refreshToken:     "refresh_token",
		expectedCode:     200,
		jwtToken:         "jwt_token",
		newRefreshToken:  "new_refresh",
		cookieExpiry:     time.Now().Add(24 * time.Hour),
		expectedResponse: model.SuccessfulAuthentication{Token: "jwt_token"},
	},
	{
		name:              utils.TestName("empty refresh token"),
		expectedCode:      400,
		emptyRefreshToken: true,
		expectedResponse:  CommonResponse{Message: "Empty refresh token. Please sign-in", Code: 400},
	},
	{
		name:             utils.TestName("active refresh token not found"),
		refreshToken:     "refresh_token",
		expectedCode:     400,
		serviceError:     model.TokenNotFound,
		expectedResponse: CommonResponse{Message: "Active refresh token not found. Please sign-in", Code: 400},
	},
	{
		name:             utils.TestName("token is expired"),
		refreshToken:     "refresh_token",
		expectedCode:     400,
		serviceError:     model.TokenExpired,
		expectedResponse: CommonResponse{Message: "Active refresh token not found. Please sign-in", Code: 400},
	},
	{
		name:             utils.TestName("failed to refresh token"),
		refreshToken:     "refresh_token",
		expectedCode:     500,
		serviceError:     errors.New("failed to refresh token"),
		expectedResponse: CommonResponse{Message: "Failed to refresh token. Please sign-in", Code: 500},
	},
}
