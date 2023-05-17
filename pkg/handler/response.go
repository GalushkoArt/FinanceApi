package handler

import (
	"FinanceApi/pkg/log"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Warn(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
