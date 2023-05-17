package handler

import (
	"FinanceApi/pkg/model"
	"FinanceApi/pkg/service"
	"github.com/GalushkoArt/simpleCache"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	sh symbolHandler
}

func New(symbolService service.SymbolService, symbolCache simpleCache.GenericCache[model.Symbol]) *Handler {
	return &Handler{sh: symbolHandler{
		service: symbolService,
		cache:   symbolCache,
	}}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			symbols := v1.Group("/symbols")
			{
				symbols.GET("", h.sh.getSymbols)
				symbols.POST("", h.sh.addSymbol)
				symbols.PUT("", h.sh.updateSymbol)
				symbols.GET("/:symbol", h.sh.getSymbol)
				symbols.DELETE("/:symbol", h.sh.deleteSymbol)
			}
		}
	}
	return router
}
