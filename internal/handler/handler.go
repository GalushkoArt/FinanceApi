package handler

import (
	"github.com/GalushkoArt/simpleCache"
	_ "github.com/galushkoart/finance-api/docs"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	swaggerHandler fiber.Handler
	ah             authHandler
	sh             symbolHandler
	apiMiddleware  []fiber.Handler
}

func New(
	swaggerHandler fiber.Handler,
	authService service.AuthService,
	symbolService service.SymbolService,
	symbolCache simpleCache.GenericCache[model.Symbol],
	apiMiddleware ...fiber.Handler,
) *Handler {
	ahLog = log.With().Str("from", "authHandler").Logger()
	shLog = log.With().Str("from", "symbolHandler").Logger()
	return &Handler{
		swaggerHandler: swaggerHandler,
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
	app.Get("/swagger/*", h.swaggerHandler)
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
				symbols.Post("", AdminOnly, h.sh.AddSymbol)
				symbols.Put("", AdminOnly, h.sh.UpdateSymbol)
				symbols.Get("/:symbol", h.sh.GetSymbol)
				symbols.Delete("/:symbol", AdminOnly, h.sh.DeleteSymbol)
			}
		}
	}
}

func setupFiberTest(handler *Handler, middleware ...func(c *fiber.Ctx) error) *fiber.App {
	app := fiber.New()
	handler.apiMiddleware = middleware
	handler.InitRoutes(app)
	return app
}
