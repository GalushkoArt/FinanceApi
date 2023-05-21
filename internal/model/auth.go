package model

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

type SignUp struct {
	Username string `json:"username" validate:"min=3,max=32,excludesall=!@#?." binding:"required"`
	Email    string `json:"email" validate:"email,max=255" binding:"required"`
	Password string `json:"password" validate:"min=6,max=32" binding:"required"`
}

type SignIn struct {
	Login    string `json:"login" validate:"min=3,max=255" binding:"required"`
	Password string `json:"password" validate:"min=6,max=32" binding:"required"`
}

type User struct {
	ID       string `json:"-"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     Role   `json:"role"`
	Password string `json:"-"`
}

type SuccessfulAuthentication struct {
	Token string `json:"token"`
}

type Role string

const (
	AdminRole  Role = "FinanceAdminRole"
	ClientRole Role = "Client"
)

var UserNotFound = errors.New("user not found")

type AuthError struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
}

var validate = validator.New()

func Validate[T SignIn | SignUp](action T) []*AuthError {
	var authErrors []*AuthError
	err := validate.Struct(action)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element AuthError
			element.Field = err.StructNamespace()
			element.Rule = err.Tag()
			authErrors = append(authErrors, &element)
		}
	}
	return authErrors
}
