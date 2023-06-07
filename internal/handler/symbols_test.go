package handler

import (
	"errors"
	"fmt"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/mock"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/golang/mock/gomock"
	"strings"
	"testing"
)

//go:generate echo $PWD - $GOFILE
//go:generate mockgen -package mock -destination ../../mock/symbol_service_mock.go -source=../service/symbol_service.go SymbolService
//use this generate function if generic cache is changed but need some manual correction of mock go:generate mockgen -package mock -destination ../../mock/generic_cache_mock.go github.com/GalushkoArt/simpleCache GenericCache[any]

func TestGetSymbols(t *testing.T) {
	controller := gomock.NewController(t)
	mockService := mock.NewMockSymbolService(controller)
	app := setupFiberTest(&Handler{sh: symbolHandler{service: mockService}}, utils.TestAuthMiddleware)
	for _, td := range getSymbolsTests {
		t.Run(td.name, func(t *testing.T) {
			mockService.EXPECT().GetAll(gomock.Any()).Return(td.symbols, td.serviceError)
			response, err := app.Test(utils.GetRequest("/api/v1/symbols"))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var getSymbolsTests = []struct {
	name             string
	symbols          []model.Symbol
	serviceError     error
	expectedCode     int
	expectedResponse interface{}
}{
	{
		name:             utils.TestName("get symbols successfully"),
		symbols:          []model.Symbol{{Symbol: "TEST"}},
		expectedCode:     200,
		expectedResponse: []model.Symbol{{Symbol: "TEST"}},
	},
	{
		name:             utils.TestName("get symbols failed"),
		serviceError:     errors.New("failed to get symbols"),
		expectedCode:     404,
		expectedResponse: CommonResponse{Code: 404, Message: "CommonResponse on retrieving all symbols"},
	},
}

// GetSymbol godoc
//
//	@Summary		GetSymbol
//	@Tags			Symbols
//	@Description	Get latest data for particular symbol
//	@Security		ApiKeyAuth[client, admin]
//	@ID				get-symbol
//	@Produce		json
//	@Success		200		{array}		model.Symbol	"Successful response"
//	@Failure		400,404	{object}	CommonResponse	"Client request error"
//	@Failure		401		{object}	CommonResponse	"Unauthorized"
//	@Failure		500		{object}	CommonResponse	"Internal server error"
//	@Router			/api/v1/symbols/{symbol} [get]
//func (h *symbolHandler) GetSymbol(c *fiber.Ctx) error {
//	symbol := strings.Replace(c.Params("symbol"), "-", "/", 1)
//	if len(symbol) == 0 {
//		return h.infoErrorResponse(c, errors.New("symbol is empty"), fiber.StatusBadRequest, "Symbol must not be empty")
//	}
//	cached := h.cache.Get(symbol)
//	if cached != nil {
//		shLog.Debug().Msgf("Return %s symbol from cache", symbol)
//		return c.Status(fiber.StatusOK).JSON(cached)
//	}
//	found, err := h.service.GetBySymbol(c.Context(), symbol)
//	if err == model.SymbolNotFound {
//		return h.infoErrorResponse(c, errors.New("symbol not found"), fiber.StatusNotFound, fmt.Sprintf("symbol %s not found", symbol))
//	} else if err != nil {
//		return h.errorErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get %s symbol", symbol))
//	}
//	h.cache.Set(symbol, found)
//	return c.Status(fiber.StatusOK).JSON(found)
//}

func TestGetSymbol(t *testing.T) {
	controller := gomock.NewController(t)
	mockService := mock.NewMockSymbolService(controller)
	mockCache := mock.NewMockGenericCache[model.Symbol](controller)
	app := setupFiberTest(&Handler{sh: symbolHandler{service: mockService, cache: mockCache}})
	for _, td := range getSymbolTests {
		t.Run(td.name, func(t *testing.T) {
			symbolParam := strings.Replace(td.requestedSymbol, "-", "/", 1)
			if td.cached {
				mockCache.EXPECT().Get(symbolParam).Return(&td.symbol)
			} else {
				mockCache.EXPECT().Get(symbolParam).Return(nil)
				mockService.EXPECT().GetBySymbol(gomock.Any(), symbolParam).Return(td.symbol, td.serviceError)
				if td.serviceError == nil {
					mockCache.EXPECT().Set(symbolParam, td.symbol).Times(1)
				}
			}
			response, err := app.Test(utils.GetRequest("/api/v1/symbols/" + td.requestedSymbol))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var getSymbolTests = []struct {
	name             string
	requestedSymbol  string
	symbol           model.Symbol
	serviceError     error
	cached           bool
	expectedCode     int
	expectedResponse interface{}
}{
	{
		name:             utils.TestName("get symbol successfully"),
		requestedSymbol:  "TE-ST",
		symbol:           model.Symbol{Symbol: "TE/ST", Name: "Test Symbol"},
		expectedCode:     200,
		expectedResponse: model.Symbol{Symbol: "TE/ST", Name: "Test Symbol"},
	},
	{
		name:             utils.TestName("get symbol from cache"),
		requestedSymbol:  "TE-ST",
		symbol:           model.Symbol{Symbol: "TE/ST"},
		cached:           true,
		expectedCode:     200,
		expectedResponse: model.Symbol{Symbol: "TE/ST"},
	},
	{
		name:             utils.TestName("symbol not found"),
		requestedSymbol:  "TE-ST",
		serviceError:     model.SymbolNotFound,
		expectedCode:     404,
		expectedResponse: CommonResponse{Code: 404, Message: "symbol TE/ST not found"},
	},
	{
		name:             utils.TestName("get symbol with internal server error"),
		requestedSymbol:  "TE-ST",
		serviceError:     errors.New("failed to get TE-ST symbol"),
		expectedCode:     500,
		expectedResponse: CommonResponse{Code: 500, Message: "Failed to get TE/ST symbol"},
	},
}

func TestDeleteSymbol(t *testing.T) {
	controller := gomock.NewController(t)
	mockService := mock.NewMockSymbolService(controller)
	mockCache := mock.NewMockGenericCache[model.Symbol](controller)
	app := setupFiberTest(&Handler{sh: symbolHandler{service: mockService, cache: mockCache}}, utils.TestAuthMiddleware)
	for _, td := range deleteSymbolTests {
		t.Run(td.name, func(t *testing.T) {
			if !td.wrongContentType && td.role == model.AdminRole {
				mockService.EXPECT().Delete(gomock.Any(), td.symbol).Return(td.serviceError)
			}
			if td.expectedCode == 200 {
				mockCache.EXPECT().Delete(td.symbol).Return(nil)
			}
			request := utils.DeleteRequest(fmt.Sprintf("/api/v1/symbols/%s", td.symbol), nil, td.wrongContentType, map[string]string{"Role": string(td.role)})
			response, err := app.Test(request)
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var deleteSymbolTests = []struct {
	name             string
	role             model.Role
	symbol           string
	expectedCode     int
	serviceError     error
	wrongContentType bool
	expectedResponse CommonResponse
}{
	{
		name:             utils.TestName("delete symbol successfully"),
		role:             model.AdminRole,
		symbol:           "TEST",
		expectedCode:     200,
		expectedResponse: CommonResponse{Code: 200, Message: "successful"},
	},
	{
		name:             utils.TestName("delete symbol with client role"),
		role:             model.ClientRole,
		symbol:           "TEST",
		expectedCode:     401,
		expectedResponse: CommonResponse{Code: 401, Message: "you don't have permissions for this endpoint"},
	},
	{
		name:             utils.TestName("delete symbol failed with not found error"),
		role:             model.AdminRole,
		symbol:           "INVALID",
		serviceError:     model.SymbolNotFound,
		expectedCode:     404,
		expectedResponse: CommonResponse{Code: 404, Message: "symbol not found"},
	},
	{
		name:             utils.TestName("delete symbol failed with service error"),
		role:             model.AdminRole,
		symbol:           "TEST",
		serviceError:     errors.New("failed to delete TEST symbol"),
		expectedCode:     500,
		expectedResponse: CommonResponse{Code: 500, Message: "Failed to delete TEST symbol"},
	},
}

func TestAddSymbol(t *testing.T) {
	mockService := mock.NewMockSymbolService(gomock.NewController(t))
	app := setupFiberTest(&Handler{sh: symbolHandler{service: mockService}}, utils.TestAuthMiddleware)
	for _, td := range addSymbolTests {
		t.Run(utils.TestName(td.name), func(t *testing.T) {
			if !td.wrongContentType && td.role == model.AdminRole {
				mockService.EXPECT().Add(gomock.Any(), td.symbol).Return(td.serviceError)
			}
			response, err := app.Test(utils.PostRequest("/api/v1/symbols", td.symbol, td.wrongContentType, map[string]string{"Role": string(td.role)}))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var addSymbolTests = []struct {
	name             string
	role             model.Role
	symbol           model.Symbol
	expectedCode     int
	serviceError     error
	wrongContentType bool
	expectedResponse CommonResponse
}{
	{
		name:             utils.TestName("add symbol successfully"),
		role:             model.AdminRole,
		symbol:           model.Symbol{Symbol: "TEST"},
		expectedCode:     200,
		expectedResponse: CommonResponse{Code: 200, Message: "successful"},
	},
	{
		name:             utils.TestName("add symbol with client role"),
		role:             model.ClientRole,
		symbol:           model.Symbol{Symbol: "TEST"},
		expectedCode:     401,
		expectedResponse: CommonResponse{Code: 401, Message: "you don't have permissions for this endpoint"},
	},
	{
		name:             utils.TestName("add symbol with wrong content type"),
		role:             model.AdminRole,
		symbol:           model.Symbol{Symbol: "TEST"},
		expectedCode:     400,
		wrongContentType: true,
		expectedResponse: CommonResponse{Code: 400, Message: "Wrong content type"},
	},
	{
		name:             utils.TestName("add symbol failed"),
		role:             model.AdminRole,
		symbol:           model.Symbol{Symbol: "TEST"},
		expectedCode:     500,
		serviceError:     errors.New("failed to add TEST symbol"),
		expectedResponse: CommonResponse{Code: 500, Message: "Failed to add TEST symbol"},
	},
}

func TestUpdateSymbol(t *testing.T) {
	controller := gomock.NewController(t)
	mockService := mock.NewMockSymbolService(controller)
	mockCache := mock.NewMockGenericCache[model.Symbol](controller)
	app := setupFiberTest(&Handler{sh: symbolHandler{service: mockService, cache: mockCache}}, utils.TestAuthMiddleware)
	for _, td := range updateSymbolTests {
		t.Run(td.name, func(t *testing.T) {
			if !td.wrongContentType && td.role == model.AdminRole {
				mockService.EXPECT().Update(gomock.Any(), td.symbol).Return(td.serviceError)
			}
			if td.expectedCode == 200 {
				mockCache.EXPECT().Delete(td.symbol.Symbol).Return(nil)
			}
			response, err := app.Test(utils.PutRequest("/api/v1/symbols", td.symbol, td.wrongContentType, map[string]string{"Role": string(td.role)}))
			utils.CommonResponseAssertions(t, response, err, td.expectedCode, td.expectedResponse)
		})
	}
}

var updateSymbolTests = []struct {
	name             string
	role             model.Role
	symbol           model.UpdateSymbol
	expectedCode     int
	serviceError     error
	wrongContentType bool
	expectedResponse CommonResponse
}{
	{
		name:             utils.TestName("update symbol successfully"),
		role:             model.AdminRole,
		symbol:           model.UpdateSymbol{Symbol: "TEST"},
		expectedCode:     200,
		expectedResponse: CommonResponse{Code: 200, Message: "successful"},
	},
	{
		name:             utils.TestName("update symbol with client role"),
		role:             model.ClientRole,
		symbol:           model.UpdateSymbol{Symbol: "TEST"},
		expectedCode:     401,
		expectedResponse: CommonResponse{Code: 401, Message: "you don't have permissions for this endpoint"},
	},
	{
		name:             utils.TestName("update symbol with wrong content type"),
		role:             model.AdminRole,
		symbol:           model.UpdateSymbol{Symbol: "TEST"},
		expectedCode:     400,
		wrongContentType: true,
		expectedResponse: CommonResponse{Code: 400, Message: "Wrong content type"},
	},
	{
		name:             utils.TestName("update symbol failed"),
		role:             model.AdminRole,
		symbol:           model.UpdateSymbol{Symbol: "TEST"},
		expectedCode:     500,
		serviceError:     errors.New("failed to update TEST symbol"),
		expectedResponse: CommonResponse{Code: 500, Message: "Failed to update TEST symbol"},
	},
}
