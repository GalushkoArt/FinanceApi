package handler

import (
	"FinanceApi/pkg/model"
	"FinanceApi/pkg/repository"
	"FinanceApi/pkg/service"
	"fmt"
	"github.com/GalushkoArt/simpleCache"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type symbolHandler struct {
	service service.SymbolService
	cache   simpleCache.GenericCache[model.Symbol]
}

func (h *symbolHandler) getSymbols(c *gin.Context) {
	symbols, err := h.service.GetAll()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, symbols)
}

func (h *symbolHandler) getSymbol(c *gin.Context) {
	symbol := strings.Replace(c.Param("symbol"), "-", "/", 1)
	if len(symbol) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "symbol must not be empty")
		return
	}
	cached := h.cache.Get(symbol)
	if cached != nil {
		c.JSON(http.StatusOK, cached)
		return
	}
	found, err := h.service.GetBySymbol(symbol)
	if err == repository.SymbolNotFound {
		newErrorResponse(c, http.StatusNotFound, fmt.Sprintf("symbol %s not found", symbol))
		return
	} else if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.cache.Set(symbol, found)
	c.JSON(http.StatusOK, found)
}

func (h *symbolHandler) deleteSymbol(c *gin.Context) {
	symbol := strings.Replace(c.Param("symbol"), "-", "/", 1)
	if len(symbol) == 0 {
		newErrorResponse(c, http.StatusBadRequest, "symbol must not be empty")
		return
	}
	if err := h.service.Delete(symbol); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.cache.Delete(symbol)
	c.Status(http.StatusNoContent)
}

func (h *symbolHandler) addSymbol(c *gin.Context) {
	var symbol model.Symbol
	err := c.BindJSON(&symbol)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.Add(symbol); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *symbolHandler) updateSymbol(c *gin.Context) {
	var symbol model.Symbol
	err := c.BindJSON(&symbol)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err = h.service.Update(symbol); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	h.cache.Delete(symbol.Symbol)
	c.Status(http.StatusNoContent)
}
