package handler

import (
	"FinanceApi/internal/model"
	"FinanceApi/internal/service"
	"errors"
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"strings"
)

type symbolHandler struct {
	service service.SymbolService
	cache   simpleCache.GenericCache[model.Symbol]
}

var shLog zerolog.Logger

// @Summary GetSymbols
// @Tag Symbols
// @Description Get all available latest symbols
// @ID get-symbols
// @Produce json
// @Success 200 {array} model.Symbol "Successful response"
// @Failure 404 {object} Error "Data not found"
// @Router /api/v1/symbols [get]
func (h *symbolHandler) getSymbols(c *fiber.Ctx) error {
	symbols, err := h.service.GetAll(c.Context())
	if err != nil {
		return warnErrorResponse(c, err, fiber.StatusNotFound, "Error on retrieving all books")
	}
	return c.Status(fiber.StatusOK).JSON(symbols)
}

// @Summary GetSymbol
// @Tag Symbols
// @Description Get latest data for particular symbol
// @ID get-symbol
// @Produce json
// @Success 200 {array} model.Symbol "Successful response"
// @Failure 400,404 {object} Error "Client request error"
// @Failure 500 {object} Error "Internal server error"
// @Router /api/v1/symbols/{symbol} [get]
func (h *symbolHandler) getSymbol(c *fiber.Ctx) error {
	symbol := strings.Replace(c.Params("symbol"), "-", "/", 1)
	if len(symbol) == 0 {
		return infoErrorResponse(c, errors.New("symbol is empty"), fiber.StatusBadRequest, "Symbol must not be empty")
	}
	cached := h.cache.Get(symbol)
	if cached != nil {
		shLog.Debug().Msgf("Return %s symbol from cache", symbol)
		return c.Status(fiber.StatusOK).JSON(cached)
	}
	found, err := h.service.GetBySymbol(c.Context(), symbol)
	if err == model.SymbolNotFound {
		return infoErrorResponse(c, errors.New("symbol not found"), fiber.StatusNotFound, fmt.Sprintf("symbol %s not found", symbol))
	} else if err != nil {
		return warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to get %s symbol", symbol))
	}
	h.cache.Set(symbol, found)
	return c.Status(fiber.StatusOK).JSON(found)
}

// @Summary DeleteSymbol
// @Tag Symbols
// @Description Delete data for symbol
// @ID delete-symbol
// @Produce json
// @Success 204 "Deleted successfully"
// @Failure 400,404 {object} Error "Client request errors"
// @Failure 500 {object} Error "Internal server errors"
// @Router /api/v1/symbols/{symbol} [delete]
func (h *symbolHandler) deleteSymbol(c *fiber.Ctx) error {
	symbol := strings.Replace(c.Params("symbol"), "-", "/", 1)
	if len(symbol) == 0 {
		return infoErrorResponse(c, errors.New("symbol is empty"), fiber.StatusBadRequest, "Symbol must not be empty")
	}
	if err := h.service.Delete(c.Context(), symbol); err != nil {
		return warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to delete %s symbol", symbol))
	}
	h.cache.Delete(symbol)
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary AddSymbols
// @Tag Symbols
// @Description Add new symbol data
// @ID add-symbols
// @Accept json
// @Param input body model.Symbol true "New symbol data"
// @Success 204 "Add successfully"
// @Failure 400 {object} Error "Client request errors"
// @Failure 500 {object} Error "Internal server errors"
// @Router /api/v1/symbols [post]
func (h *symbolHandler) addSymbol(c *fiber.Ctx) error {
	var symbol model.Symbol
	if err := c.BodyParser(&symbol); err != nil {
		return infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong body")
	}
	if err := h.service.Add(c.Context(), symbol); err != nil {
		return warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to add %s symbol", symbol.Symbol))
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary UpdateSymbols
// @Tag Symbols
// @Description Update symbol data
// @ID update-symbols
// @Accept json
// @Param input body model.UpdateSymbol true "Update symbol data"
// @Success 204 "Add successfully"
// @Failure 400 {object} Error "Client request errors"
// @Failure 500 {object} Error "Internal server errors"
// @Router /api/v1/symbols [put]
func (h *symbolHandler) updateSymbol(c *fiber.Ctx) error {
	var symbol model.UpdateSymbol
	if err := c.BodyParser(&symbol); err != nil {
		return infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong body")
	}
	if err := h.service.Update(c.Context(), symbol); err != nil {
		return warnErrorResponse(c, err, fiber.StatusInternalServerError, fmt.Sprintf("Failed to update %s symbol", symbol.Symbol))
	}
	h.cache.Delete(symbol.Symbol)
	return c.SendStatus(fiber.StatusNoContent)
}
