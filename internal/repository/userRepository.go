package repository

import (
	"FinanceApi/internal/model"
	"FinanceApi/pkg/utils"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type userRepositoryPostgres struct {
	db *sqlx.DB
}

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	CheckLoginIsAvailable(ctx context.Context, username string, email string) (bool, error)
	GetUser(ctx context.Context, login string, password string) (model.User, error)
	GetUserRole(ctx context.Context, userID string) (model.Role, error)
	Update(ctx context.Context, user model.User) error
	GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error)
	InsertRefreshToken(ctx context.Context, token model.RefreshToken) error
}

func urLog(c context.Context, e *zerolog.Event) *zerolog.Event {
	return utils.LogRequest(c, e).Str("from", "userRepositoryPostgres")
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepositoryPostgres{db: db}
}

const logInsert = `INSERT INTO LOGS(EVENT) VALUES ($1)`

func (r *userRepositoryPostgres) Create(ctx context.Context, user model.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Failed to begin transaction")
		return err
	}
	urLog(ctx, log.Info()).Msgf("Creating user: %s username, %s email", user.Username, user.Email)
	const userInsert = `INSERT INTO USER_ENTITY(ID, USERNAME, EMAIL, PASSWORD) VALUES ($1, $2, $3, $4) RETURNING id`
	_, err = tx.Exec(userInsert, user.ID, user.Username, user.Email, user.Password)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on insert user!")
		utils.PanicOnError(tx.Rollback())
		return err
	}
	urLog(ctx, log.Debug()).Msgf("Add %s role for %s user", user.Username, user.Email)
	const userRoleInsert = `INSERT INTO user_roles(user_id, role) VALUES ($1, $2)`
	_, err = tx.Exec(userRoleInsert, user.ID, user.Role)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on insert user role!")
		utils.PanicOnError(tx.Rollback())
		return err
	}
	urLog(ctx, log.Debug()).Msg("User created! Write logs to DB")
	_, err = tx.Exec(logInsert, fmt.Sprintf("Created user with %s id, %s username, %s email", user.ID, user.Username, user.Email))
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on insert logs!")
		utils.PanicOnError(tx.Rollback())
		return err
	}
	return tx.Commit()
}

func (r *userRepositoryPostgres) CheckLoginIsAvailable(ctx context.Context, username string, email string) (bool, error) {
	var count int
	urLog(ctx, log.Debug()).Msgf("Checking are %s username and %s email available", username, email)
	const countUsernameEmailEntriesQuery = `SELECT COUNT(1) FROM user_entity WHERE username = $1 or email = $2`
	err := r.db.GetContext(ctx, &count, countUsernameEmailEntriesQuery, username, email)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *userRepositoryPostgres) GetUser(ctx context.Context, login string, password string) (model.User, error) {
	var users []userEntity
	urLog(ctx, log.Debug()).Msgf("Retrieving user with %s login", login)
	const usersByLoginQuery = `SELECT * FROM user_entity WHERE (username = $1 or email = $1) and password = $2`
	err := r.db.SelectContext(ctx, &users, usersByLoginQuery, login, password)
	if err != nil {
		return model.User{}, err
	}
	if len(users) == 0 {
		return model.User{}, model.UserNotFound
	}
	role, err := r.GetUserRole(ctx, users[0].ID)
	if len(users) == 0 {
		return model.User{}, err
	}
	result := model.User{
		ID:       users[0].ID,
		Username: users[0].Username,
		Email:    users[0].Email,
		Role:     role,
		Password: users[0].Password,
	}
	return result, nil
}
func (r *userRepositoryPostgres) GetUserRole(ctx context.Context, userID string) (model.Role, error) {
	var role string
	urLog(ctx, log.Debug()).Msgf("Retrieving user role with %s user id", userID)
	const usersRoleQuery = `SELECT role FROM user_roles WHERE user_id = $1`
	err := r.db.GetContext(ctx, &role, usersRoleQuery, userID)
	if err != nil {
		return model.ClientRole, err
	}
	return model.Role(role), nil
}

func (r *userRepositoryPostgres) Update(ctx context.Context, user model.User) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Failed to begin transaction")
		return err
	}
	urLog(ctx, log.Info()).Msgf("Updating user: %s username, %s email", user.Username, user.Email)
	const userUpdate = `UPDATE USER_ENTITY SET USERNAME = $1, EMAIL = $2, PASSWORD = $3, UPDATED_AT = now() WHERE ID = $4`
	_, err = tx.Exec(userUpdate, user.Username, user.Email, user.Password, user.ID)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on update user!")
		utils.PanicOnError(tx.Rollback())
		return err
	}
	urLog(ctx, log.Debug()).Msg("User updated! Write logs to DB")
	_, err = tx.Exec(logInsert, fmt.Sprintf("Update user with %s id, %s username, %s email", user.ID, user.Username, user.Email))
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on insert logs!")
		utils.PanicOnError(tx.Rollback())
		return err
	}
	return tx.Commit()
}

func (r *userRepositoryPostgres) GetRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Failed to begin transaction")
		return model.RefreshToken{}, err
	}
	var tokens []refreshToken
	urLog(ctx, log.Debug()).Msgf("Retrieving refresh token with %s value", token)
	const usersByLoginQuery = `SELECT * FROM refresh_token WHERE token = $1`
	err = tx.SelectContext(ctx, &tokens, usersByLoginQuery, token)
	if err != nil {
		return model.RefreshToken{}, err
	}
	if len(tokens) == 0 {
		return model.RefreshToken{}, model.TokenNotFound
	}
	urLog(ctx, log.Debug()).Msgf("Deleting old token: %s", token)
	const deleteToken = `DELETE FROM refresh_token WHERE ID = $1`
	_, err = tx.Exec(deleteToken, tokens[0].ID)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail to delete token")
		utils.PanicOnError(tx.Rollback())
		return model.RefreshToken{}, err
	}
	result := model.RefreshToken{
		Token:     token,
		UserId:    tokens[0].UserId,
		ExpiresAt: tokens[0].ExpiresAt,
	}
	return result, nil
}

func (r *userRepositoryPostgres) InsertRefreshToken(ctx context.Context, refreshToken model.RefreshToken) error {
	urLog(ctx, log.Info()).Msgf("Inserting refresh token for %s user id", refreshToken.UserId)
	const insertRefreshToken = `INSERT INTO refresh_token(USER_ID, TOKEN, EXPIRES_AT) VALUES ($1, $2, $3)`
	_, err := r.db.ExecContext(ctx, insertRefreshToken, refreshToken.UserId, refreshToken.Token, refreshToken.ExpiresAt)
	if err != nil {
		urLog(ctx, log.Error()).Err(err).Msg("Fail on insert token!")
	}
	return err
}
