package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
)

type UserService struct {
	Pool    *pgxpool.Pool
	Queries *dbqueries.Queries
}

func NewUserService(pool *pgxpool.Pool, queries *dbqueries.Queries) *UserService {
	return &UserService{
		Pool:    pool,
		Queries: queries,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username string, passwd string, role int32, deviceName string) (dbqueries.User, dbqueries.Device, error) {
	// check if user already exists
	oldUser, err := s.Queries.GetUserByName(ctx, username)
	if oldUser.UserID.String() != "" {
		return dbqueries.User{}, dbqueries.Device{}, ErrUserAlreadyExists
	}

	// create password hash
	hash, err := argon2id.CreateHash(passwd, argon2id.DefaultParams)
	if err != nil {
		return dbqueries.User{}, dbqueries.Device{}, err
	}
	userParams := dbqueries.CreateUserParams{
		Username: username,
		Pwhash:   hash,
		Role:     role,
	}

	// create user
	user, err := s.Queries.CreateUser(ctx, userParams)
	if err != nil {
		return dbqueries.User{}, dbqueries.Device{}, err
	}

	token, err := randomHex(32)
	if err != nil {
		return dbqueries.User{}, dbqueries.Device{}, err
	}
	deviceParams := dbqueries.CreateDeviceParams{
		UserID:     user.UserID,
		LastSeen:   pgtype.Timestamptz{Time: time.Now(), Valid: true},
		LoginToken: token,
		DeviceName: deviceName,
	}

	// create device
	device, err := s.Queries.CreateDevice(ctx, deviceParams)
	if err != nil {
		return dbqueries.User{}, dbqueries.Device{}, err
	}

	return user, device, nil
}

func randomHex(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(bytes)), nil
}
