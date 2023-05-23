package handler

import (
	"FinanceApi/internal/model"
	"FinanceApi/internal/service"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"strings"
)

type authHandler struct {
	service service.AuthService
}

var ahLog zerolog.Logger

func (h *authHandler) errorErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return errorErrorResponse(c, &ahLog, err, statusCode, message, authErrors...)
}

func (h *authHandler) warnErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return warnErrorResponse(c, &ahLog, err, statusCode, message, authErrors...)
}

func (h *authHandler) infoErrorResponse(c *fiber.Ctx, err error, statusCode int, message string, authErrors ...[]*model.AuthError) error {
	return infoErrorResponse(c, &ahLog, err, statusCode, message, authErrors...)
}

const refreshTokenCookie = "refresh-token"

// SignUp godoc
//
//	@Summary		SignUp
//	@Tags			Auth
//	@Description	Register new user
//	@ID				sign-up
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.SignUp	true	"New user data"
//	@Success		200		{object}	CommonResponse	"New user created successfully"
//	@Failure		400		{object}	CommonResponse	"Wrong user data"
//	@Failure		500		{object}	CommonResponse	"Internal server errors"
//	@Router			/auth/signup [put]
func (h *authHandler) SignUp(c *fiber.Ctx) error {
	var signUp model.SignUp
	if err := c.BodyParser(&signUp); err != nil {
		return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong body")
	}
	signUp.Username = strings.ToLower(signUp.Username)
	signUp.Email = strings.ToLower(signUp.Email)
	authErrors := model.Validate(signUp)
	if authErrors != nil {
		return h.infoErrorResponse(c, errors.New("invalid sign-up body"), fiber.StatusBadRequest, "Wrong body", authErrors)
	}
	if err := h.service.SignUp(c.Context(), signUp); err != nil {
		if err == service.UserAlreadyExists {
			return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "User with such username or email already exists")
		}
		return h.warnErrorResponse(c, err, fiber.StatusInternalServerError, "Failed to register")
	}
	return c.Status(fiber.StatusOK).JSON(CommonResponse{Code: fiber.StatusOK, Message: "success"})
}

// SignIn godoc
//
//	@Summary		SignIn
//	@Tags			Auth
//	@Description	Authenticate user
//	@ID				sign-in
//	@Accept			json
//	@Produce		json
//	@Param			input	body		model.SignIn					true	"Authentication user data"
//	@Success		200		{object}	model.SuccessfulAuthentication	"Response with jwt token"
//	@Failure		400		{object}	CommonResponse					"Wrong user data"
//	@Failure		401		{object}	CommonResponse					"Wrong credentials"
//	@Failure		500		{object}	CommonResponse					"Internal server errors"
//	@Router			/auth/signin [put]
func (h *authHandler) SignIn(c *fiber.Ctx) error {
	var signIn model.SignIn
	if err := c.BodyParser(&signIn); err != nil {
		return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "Wrong body")
	}
	signIn.Login = strings.ToLower(signIn.Login)
	authErrors := model.Validate(signIn)
	if authErrors != nil {
		return h.infoErrorResponse(c, errors.New("invalid sign-in body"), fiber.StatusBadRequest, "Wrong body", authErrors)
	}
	jwtToken, refreshToken, expiryTime, err := h.service.SignIn(c.Context(), signIn)
	if err != nil {
		if err == model.UserNotFound {
			return h.infoErrorResponse(c, errors.New("wrong credentials"), fiber.StatusUnauthorized, "Wrong credentials", authErrors)
		}
		return h.errorErrorResponse(c, err, fiber.StatusInternalServerError, "Failed to sign in")
	}
	c.Cookie(&fiber.Cookie{Name: refreshTokenCookie, Value: refreshToken, Expires: expiryTime, HTTPOnly: true})
	return c.Status(fiber.StatusOK).JSON(model.SuccessfulAuthentication{Token: jwtToken})
}

// Refresh godoc
//
//	@Summary		Refresh
//	@Tags			Auth
//	@Description	Refresh auth token
//	@ID				refresh-token
//	@Produce		json
//	@Success		200	{object}	model.SuccessfulAuthentication	"Response with jwt token"
//	@Failure		400	{object}	CommonResponse					"Wrong refresh token"
//	@Failure		500	{object}	CommonResponse					"Internal server errors"
//	@Router			/auth/refresh [get]
func (h *authHandler) Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies(refreshTokenCookie, "")
	if len(refreshToken) == 0 {
		return h.infoErrorResponse(c, errors.New("empty refresh token"), fiber.StatusBadRequest, "Empty refresh token. Please sign-in")
	}
	jwtToken, refreshToken, expiryTime, err := h.service.RefreshToken(c.Context(), refreshToken)
	if err != nil {
		if err == model.TokenExpired || err == model.TokenNotFound {
			return h.infoErrorResponse(c, err, fiber.StatusBadRequest, "Active refresh token not found. Please sign-in")
		}
		return h.errorErrorResponse(c, err, fiber.StatusInternalServerError, "Failed to refresh token. Please sign-in")
	}
	c.Cookie(&fiber.Cookie{Name: refreshTokenCookie, Value: refreshToken, Expires: expiryTime, HTTPOnly: true})
	return c.Status(fiber.StatusOK).JSON(model.SuccessfulAuthentication{Token: jwtToken})
}
