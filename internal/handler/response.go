package handler

import (
	"github.com/galushkoart/finance-api/internal/model"
	"github.com/galushkoart/finance-api/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type CommonResponse struct {
	Code       int                `json:"code"`
	Message    string             `json:"message"`
	AuthErrors []*model.AuthError `json:"authErrors,omitempty"`
}

func errorErrorResponse(c *fiber.Ctx, logger *zerolog.Logger, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	utils.LogRequest(c.Context(), logger.Error()).Stack().Err(err).Msg(message)
	return returnError(c, statusCode, message, authErrors...)
}

func warnErrorResponse(c *fiber.Ctx, logger *zerolog.Logger, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	utils.LogRequest(c.Context(), logger.Warn()).Stack().Err(err).Msg(message)
	return returnError(c, statusCode, message, authErrors...)
}

func infoErrorResponse(c *fiber.Ctx, logger *zerolog.Logger, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	utils.LogRequest(c.Context(), logger.Info()).Err(err).Msg(message)
	return returnError(c, statusCode, message, authErrors...)
}

func returnError(c *fiber.Ctx, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	if len(authErrors) > 0 {
		return c.Status(statusCode).JSON(CommonResponse{Code: statusCode, Message: message, AuthErrors: authErrors[0]})
	}
	return c.Status(statusCode).JSON(CommonResponse{Code: statusCode, Message: message})
}
