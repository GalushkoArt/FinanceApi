package handler

import (
	"FinanceApi/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func warnErrorResponse(c *fiber.Ctx, err error, statusCode int, message string) error {
	shLog.Warn().Stack().Str("request-id", utils.GetRequestId(c.Context())).Err(err).Msg(message)
	return returnError(c, statusCode, message)
}

func infoErrorResponse(c *fiber.Ctx, err error, statusCode int, message string) error {
	shLog.Info().Stack().Str("request-id", utils.GetRequestId(c.Context())).Err(err).Msg(message)
	return returnError(c, statusCode, message)
}

func returnError(c *fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Error{Code: statusCode, Message: message})
}
