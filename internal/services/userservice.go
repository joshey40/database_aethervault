package services

import (
	"context"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
	"github.com/joshey40/database_aethervault/internal/logger"
	"go.uber.org/zap"
)

type UserService struct {
	Pool      *pgxpool.Pool
	Queries   *dbqueries.Queries
	TokenAuth *jwtauth.JWTAuth
}

func NewUserService(pool *pgxpool.Pool, queries *dbqueries.Queries, tokenAuth *jwtauth.JWTAuth) *UserService {
	return &UserService{
		Pool:      pool,
		Queries:   queries,
		TokenAuth: tokenAuth,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username string, passwd string, role int32, deviceName string) (string, error) {
	// check if user already exists
	oldUser, err := s.Queries.GetUserByName(ctx, username)
	if oldUser.UserID.String() != "" {
		return "", ErrUserAlreadyExists
	}

	// create password hash
	hash, err := argon2id.CreateHash(passwd, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	userParams := dbqueries.CreateUserParams{
		Username: username,
		Pwhash:   hash,
		Role:     role,
	}

	// create user
	user, err := s.Queries.CreateUser(ctx, userParams)
	if err != nil {
		return "", err
	}

	// create device
	deviceParams := dbqueries.CreateDeviceParams{
		UserID:     user.UserID,
		LastSeen:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		DeviceName: deviceName,
	}

	device, err := s.Queries.CreateDevice(ctx, deviceParams)
	if err != nil {
		return "", err
	}

	// create JWT
	_, tokenString, _ := s.TokenAuth.Encode(map[string]interface{}{"user_id": user.UserID, "device_id": device.DeviceID})

	return tokenString, nil
}

func (s *UserService) LoginUser(ctx context.Context, username string, passwd string, deviceName string) (string, error) {
	// get user
	user, err := s.Queries.GetUserByName(ctx, username)
	if err != nil {
		logger.L().Warn("login tried for nonexisting user", zap.String("username", username), zap.Error(err))
		return "", ErrInvalidLoginData
	}

	// check password
	match, err := argon2id.ComparePasswordAndHash(passwd, user.Pwhash)
	if !match {
		logger.L().Warn("login attempt with wrong password", zap.String("username", username))
		return "", ErrInvalidLoginData
	}

	// create device
	deviceParams := dbqueries.CreateDeviceParams{
		UserID:     user.UserID,
		LastSeen:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		DeviceName: deviceName,
	}

	device, err := s.Queries.CreateDevice(ctx, deviceParams)
	if err != nil {
		return "", err
	}

	// create JWT
	_, tokenString, _ := s.TokenAuth.Encode(map[string]interface{}{"user_id": user.UserID, "device_id": device.DeviceID})

	return tokenString, nil
}
