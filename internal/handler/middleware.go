package handler

import (
	"errors"
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/internal/service"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

type LogConfig struct {
	Logger *zerolog.Logger
	Next   func(c *fiber.Ctx) bool
}

func RequestLogger(logConf ...LogConfig) fiber.Handler {
	var logger zerolog.Logger
	var conf LogConfig
	if len(logConf) > 0 {
		logger = *logConf[0].Logger
		conf = logConf[0]
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

type AuthConfig struct {
	Next func(c *fiber.Ctx) bool
}

func AuthMiddleware(parser *service.JwtParser, authConf ...AuthConfig) fiber.Handler {
	var conf AuthConfig
	if len(authConf) != 0 {
		conf = authConf[0]
	}
	return func(c *fiber.Ctx) error {
		if conf.Next != nil && conf.Next(c) {
			return c.Next()
		}
		tokenString, err := getTokenFromRequest(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(CommonResponse{Code: fiber.StatusUnauthorized, Message: err.Error()})
		}
		userId, role, err := parser.ParseToken(tokenString)
		if err != nil {
			if err == jwt.ErrHashUnavailable {
				return c.SendStatus(fiber.StatusInternalServerError)
			}
			return c.Status(fiber.StatusUnauthorized).JSON(CommonResponse{Code: fiber.StatusUnauthorized, Message: err.Error()})
		}
		log.Info().
			Str("request-id", utils.GetRequestId(c.Context())).
			Str("from", "authMiddleware").
			Msgf("request from %s with %s id", role, userId)
		c.Locals("role", role)
		return c.Next()
	}
}

func getTokenFromRequest(c *fiber.Ctx) (string, error) {
	authHeader := c.GetReqHeaders()["Authorization"]
	if authHeader == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}

func AdminOnly(c *fiber.Ctx) error {
	if c.Locals("role") != model.AdminRole {
		return c.Status(fiber.StatusUnauthorized).JSON(CommonResponse{Code: fiber.StatusUnauthorized, Message: "you don't have permissions for this endpoint"})
	}
	return c.Next()
}
