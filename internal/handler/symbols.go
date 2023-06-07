package handler

import (
	"errors"
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"strings"
)

type symbolHandler struct {
	service service.SymbolService
	cache   simpleCache.GenericCache[model.Symbol]
}

var shLog zerolog.Logger

func (h *symbolHandler) errorErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return errorErrorResponse(c, &shLog, err, statusCode, message, authErrors...)
}

func (h *symbolHandler) warnErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return warnErrorResponse(c, &shLog, err, statusCode, message, authErrors...)
}

func (h *symbolHandler) infoErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return infoErrorResponse(c, &shLog, err, statusCode, message, authErrors...)
}

// GetSymbols godoc
//
//	@Summary		GetSymbols
//	@Tags			Symbols
//	@Description	Get all available latest symbols
//	@Security		ApiKeyAuth[client, admin]
//	@ID				get-symbols
//	@Produce		json
//	@Success		200	{array}		model.Symbol	"Successful response"
//	@Failure		401	{object}	CommonResponse	"Unauthorized"
//	@Failure		404	{object}	CommonResponse	"Data not found"
//	@Router			/api/v1/symbols [get]
func (h *symbolHandler) GetSymbols(c *fiber.Ctx) error {
	symbols, err := h.service.GetAll(c.Context())
	if err != nil {
		return h.warnErrorResponse(c, err, fiber.StatusNotFound, "CommonResponse on retrieving all symbols")
	}
	return c.Status(fiber.StatusOK).JSON(symbols)
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
func (h *symbolHandler) GetSymbol(c *fiber.Ctx) error {
	symbol := strings.Replace(c.Params("symbol"), "-", "/", 1)
	cached := h.cache.Get(symbol)
	if cached != nil {
		shLog.Debug().Msgf("Return %s symbol from cache", symbol)
		return c.Status(fiber.StatusOK).JSON(cached)
	}
	found, err := h.service.GetBySymbol(c.Context(), symbol)
	if err == model.SymbolNotFound {
		return h.infoErrorResponse(c, errors.New("symbol not found"), fiber.StatusNotFound, fmt.Sprintf("symbol %s not found", symbol))
	} else if err != nil {
		return h.errorErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get %s symbol", symbol))
	}
	h.cache.Set(symbol, found)
	return c.Status(fiber.StatusOK).JSON(found)
}

// DeleteSymbol godoc
//
//	@Summary		DeleteSymbol
//	@Tags			Symbols
//	@Description	Delete data for symbol
//	@Security		ApiKeyAuth[admin]
//	@ID				delete-symbol
//	@Produce		json
//	@Success		200		{object}	CommonResponse	"Deleted successfully"
//	@Failure		400,404	{object}	CommonResponse	"Client request errors"
//	@Failure		401		{object}	CommonResponse	"Unauthorized"
//	@Failure		500		{object}	CommonResponse	"Internal server errors"
//	@Router			/api/v1/symbols/{symbol} [delete]
func (h *symbolHandler) DeleteSymbol(c *fiber.Ctx) error {
	symbol := strings.Replace(c.Params("symbol"), "-", "/", 1)
	if err := h.service.Delete(c.Context(), symbol); err != nil {
		if err == model.SymbolNotFound {
			return h.infoErrorResponse(c, err, fiber.StatusNotFound, err.Error())
		}
		return h.warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to delete %s symbol", symbol))
	}
	h.cache.Delete(symbol)
	return c.Status(fiber.StatusOK).JSON(CommonResponse{Code: fiber.StatusOK, Message: "successful"})
}

// AddSymbol godoc
//
//	@Summary		AddSymbols
//	@Tags			Symbols
//	@Description	Add new symbol data
//	@Security		ApiKeyAuth[client, admin]
//	@ID				add-symbols
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.Symbol	true	"New symbol data"
//	@Success		200		{object}	CommonResponse	"Add successfully"
//	@Failure		400		{object}	CommonResponse	"Client request errors"
//	@Failure		401		{object}	CommonResponse	"Unauthorized"
//	@Failure		500		{object}	CommonResponse	"Internal server errors"
//	@Router			/api/v1/symbols [post]
func (h *symbolHandler) AddSymbol(c *fiber.Ctx) error {
	var symbol model.Symbol
	if err := c.BodyParser(&symbol); err != nil {
		return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong content type")
	}
	if err := h.service.Add(c.Context(), symbol); err != nil {
		return h.warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to add %s symbol", symbol.Symbol))
	}
	return c.Status(fiber.StatusOK).JSON(CommonResponse{Code: fiber.StatusOK, Message: "successful"})
}

// UpdateSymbol godoc
//
//	@Summary		UpdateSymbols
//	@Tags			Symbols
//	@Description	Update symbol data
//	@Security		ApiKeyAuth[admin]
//	@ID				update-symbols
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.UpdateSymbol	true	"Update symbol data"
//	@Success		200		{object}	CommonResponse		"Add successfully"
//	@Failure		400		{object}	CommonResponse		"Client request errors"
//	@Failure		401		{object}	CommonResponse		"Unauthorized"
//	@Failure		500		{object}	CommonResponse		"Internal server errors"
//	@Router			/api/v1/symbols [put]
func (h *symbolHandler) UpdateSymbol(c *fiber.Ctx) error {
	var symbol model.UpdateSymbol
	if err := c.BodyParser(&symbol); err != nil {
		return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong content type")
	}
	if err := h.service.Update(c.Context(), symbol); err != nil {
		return h.warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update %s symbol", symbol.Symbol))
	}
	h.cache.Delete(symbol.Symbol)
	return c.Status(fiber.StatusOK).JSON(CommonResponse{Code: fiber.StatusOK, Message: "successful"})
}
