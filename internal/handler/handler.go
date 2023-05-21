package handler

import (
	_ "FinanceApi/docs"
	"FinanceApi/internal/model"
	"FinanceApi/internal/service"
	"github.com/GalushkoArt/simpleCache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	sh symbolHandler
}

func New(symbolService service.SymbolService, symbolCache simpleCache.GenericCache[model.Symbol]) *Handler {
	shLog = log.With().Str("from", "symbolHandler").Logger()
	return &Handler{sh: symbolHandler{
		service: symbolService,
		cache:   symbolCache,
	}}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	api := app.Group("/api")
	{
		api.Use(RequestLogger())
		v1 := api.Group("/v1")
		{
			symbols := v1.Group("/symbols")
			{
				symbols.Get("", h.sh.getSymbols)
				symbols.Post("", h.sh.addSymbol)
				symbols.Put("", h.sh.updateSymbol)
				symbols.Get("/:symbol", h.sh.getSymbol)
				symbols.Delete("/:symbol", h.sh.deleteSymbol)
			}
		}
	}
}
