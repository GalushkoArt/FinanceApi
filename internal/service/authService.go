package service

import (
	"FinanceApi/internal/model"
	"FinanceApi/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"math/rand"
	"time"
)

type authService struct {
	repo                repository.UserRepository
	hasher              *Hasher
	jwtProducer         *JwtProducer
	refreshTokenTimeout time.Duration
	auditService        AuditService
}

type AuthService interface {
	SignUp(ctx context.Context, signUp model.SignUp) error
	SignIn(ctx context.Context, signIn model.SignIn) (string, string, time.Time, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, string, time.Time, error)
}

func NewAuthService(repo repository.UserRepository, hasher *Hasher, jwtProducer *JwtProducer, refreshTokenTimeout time.Duration, auditService AuditService) AuthService {
	return &authService{hasher: hasher, repo: repo, jwtProducer: jwtProducer, refreshTokenTimeout: refreshTokenTimeout, auditService: auditService}
}

var UserAlreadyExists = errors.New("user is already exists")

func (s *authService) SignUp(ctx context.Context, signUp model.SignUp) error {
	available, err := s.repo.CheckLoginIsAvailable(ctx, signUp.Username, signUp.Email)
	if err != nil {
		return err
	}
	if !available {
		return UserAlreadyExists
	}
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}
	go s.auditService.LogUserSignUp(ctx, id.String())
	passHash, err := s.hasher.Hash(signUp.Password)
	if err != nil {
		return err
	}
	return s.repo.Create(ctx, model.User{ID: id.String(), Username: signUp.Username, Email: signUp.Email, Role: model.ClientRole, Password: passHash})
}

func (s *authService) SignIn(ctx context.Context, signIn model.SignIn) (string, string, time.Time, error) {
	passHash, err := s.hasher.Hash(signIn.Password)
	if err != nil {
		return "", "", time.Time{}, err
	}
	user, err := s.repo.GetUser(ctx, signIn.Login, passHash)
	if err != nil {
		return "", "", time.Time{}, err
	}
	go s.auditService.LogUserSignIn(ctx, user.ID)
	return s.getTokens(ctx, user.ID, model.ClientRole)
}

func (s *authService) getTokens(ctx context.Context, userId string, role model.Role) (string, string, time.Time, error) {
	jwtToken, err := s.jwtProducer.GetToken(ctx, userId, role)
	if err != nil {
		return "", "", time.Time{}, err
	}
	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", time.Time{}, err
	}
	expiryTime := time.Now().Add(s.refreshTokenTimeout)
	err = s.repo.InsertRefreshToken(ctx, model.RefreshToken{Token: refreshToken, UserId: userId, ExpiresAt: expiryTime})
	if err != nil {
		return "", "", time.Time{}, err
	}
	return jwtToken, refreshToken, expiryTime, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (string, string, time.Time, error) {
	token, err := s.repo.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", "", time.Time{}, err
	}
	if token.ExpiresAt.Before(time.Now()) {
		return "", "", time.Time{}, model.TokenExpired
	}
	go s.auditService.LogUserRefreshToken(ctx, token.UserId)
	role, err := s.repo.GetUserRole(ctx, token.UserId)
	if err != nil {
		return "", "", time.Time{}, err
	}
	return s.getTokens(ctx, token.UserId, role)
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
