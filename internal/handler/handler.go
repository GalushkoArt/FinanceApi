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
	ah            authHandler
	sh            symbolHandler
	apiMiddleware []fiber.Handler
}

func New(
	authService service.AuthService,
	symbolService service.SymbolService,
	symbolCache simpleCache.GenericCache[model.Symbol],
	apiMiddleware ...fiber.Handler,
) *Handler {
	ahLog = log.With().Str("from", "authHandler").Logger()
	shLog = log.With().Str("from", "symbolHandler").Logger()
	return &Handler{
		ah: authHandler{
			service: authService,
		},
		sh: symbolHandler{
			service: symbolService,
			cache:   symbolCache,
		},
		apiMiddleware: apiMiddleware,
	}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
	auth := app.Group("/auth")
	{
		auth.Post("/signup", h.ah.SignUp)
		auth.Post("/signin", h.ah.SignIn)
		auth.Get("/refresh", h.ah.Refresh)
	}
	api := app.Group("/api")
	{
		for _, middleware := range h.apiMiddleware {
			api.Use(middleware)
		}
		v1 := api.Group("/v1")
		{
			symbols := v1.Group("/symbols")
			{
				symbols.Get("", h.sh.GetSymbols)
				symbols.Post("", h.sh.AddSymbol)
				symbols.Put("", h.sh.UpdateSymbol)
				symbols.Get("/:symbol", h.sh.GetSymbol)
				symbols.Delete("/:symbol", h.sh.DeleteSymbol)
			}
		}
	}
}
