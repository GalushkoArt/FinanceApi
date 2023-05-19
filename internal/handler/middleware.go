package handler

import (
	"FinanceApi/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

type Config struct {
	Logger *zerolog.Logger
	Next   func(c *fiber.Ctx) bool
}

func RequestLogger(confs ...Config) fiber.Handler {
	var logger zerolog.Logger
	var conf Config
	if len(confs) > 0 {
		logger = *confs[0].Logger
		conf = confs[0]
	} else {
		logger = log.Logger
	}
	return func(c *fiber.Ctx) error {
		if conf.Next != nil && conf.Next(c) {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()

		msg := "Request"
		if err != nil {
			msg = err.Error()
		}
		code := c.Response().StatusCode()

		resultLogger := logger.With().
			Str("request-id", utils.GetRequestId(c.Context())).
			Int("status", c.Response().StatusCode()).
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("body", string(c.Body())).
			Str("latency", time.Since(start).String()).
			Logger()

		switch {
		case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
			resultLogger.Warn().Msg(msg)
		case code >= fiber.StatusInternalServerError:
			resultLogger.Error().Msg(msg)
		default:
			resultLogger.Info().Msg(msg)
		}
		return err
	}
}
