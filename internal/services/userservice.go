package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
	"github.com/lestrrat-go/jwx/v3/jwt"
	"go.uber.org/zap"
)

type Authservice interface {
	CheckJWT(ctx context.Context) (string, string, error)
}

type UserService struct {
	DB        *sql.DB
	Queries   *dbqueries.Queries
	TokenAuth *jwtauth.JWTAuth
}

func NewUserService(db *sql.DB, queries *dbqueries.Queries, tokenAuth *jwtauth.JWTAuth) *UserService {
	return &UserService{
		DB:        db,
		Queries:   queries,
		TokenAuth: tokenAuth,
	}
}

func (s *UserService) CreateUser(ctx context.Context, username string, passwd string, role int64, deviceName string) (string, error) {
	// check if user already exists
	_, err := s.Queries.GetUserByName(ctx, username)
	if err == nil {
		return "", ErrUserAlreadyExists
	}

	// create password hash
	hash, err := argon2id.CreateHash(passwd, argon2id.DefaultParams)
	if err != nil {
		zap.L().Error("could not generate PW hash", zap.String("passwd entered", passwd), zap.Error(err))
		return "", ErrInternal
	}

	// create user uuid
	userUUID, err := uuid.NewV7()
	if err != nil {
		zap.L().Error("could not generate UUID", zap.Error(err))
		return "", ErrInternal
	}

	// create user
	userParams := dbqueries.CreateUserParams{
		UserID:   userUUID,
		Username: username,
		Pwhash:   hash,
		Role:     role,
	}
	user, err := s.Queries.CreateUser(ctx, userParams)
	if err != nil {
		zap.L().Error("could not create user", zap.Error(err))
		return "", ErrInternal
	}

	// create device uuid
	deviceUUID, err := uuid.NewV7()
	if err != nil {
		zap.L().Error("could not generate UUID", zap.Error(err))
		return "", ErrInternal
	}

	// create device
	deviceParams := dbqueries.CreateDeviceParams{
		DeviceID:   deviceUUID,
		UserID:     user.UserID,
		LastSeen:   time.Now().UnixMilli(),
		DeviceName: deviceName,
	}
	device, err := s.Queries.CreateDevice(ctx, deviceParams)
	if err != nil {
		zap.L().Error("could not create device", zap.Error(err))
		return "", ErrInternal
	}

	// create JWT
	_, tokenString, _ := s.TokenAuth.Encode(map[string]interface{}{"user_id": user.UserID, "device_id": device.DeviceID})

	return tokenString, nil
}

func (s *UserService) LoginUser(ctx context.Context, username string, passwd string, deviceName string) (string, error) {
	// get user
	user, err := s.Queries.GetUserByName(ctx, username)
	if err != nil {
		zap.L().Warn("login tried for nonexisting user", zap.String("username", username), zap.Error(err))
		return "", ErrInvalidLoginData
	}

	// check password
	match, err := argon2id.ComparePasswordAndHash(passwd, user.Pwhash)
	if !match {
		zap.L().Warn("login attempt with wrong password", zap.String("username", username))
		return "", ErrInvalidLoginData
	}
	if err != nil {
		zap.L().Error("could not compare password to passwordhash", zap.Error(err))
		return "", ErrInternal
	}

	// create device uuid
	deviceUUID, err := uuid.NewV7()
	if err != nil {
		zap.L().Error("could not generate UUID", zap.Error(err))
		return "", ErrInternal
	}

	// create device
	deviceParams := dbqueries.CreateDeviceParams{
		DeviceID:   deviceUUID,
		UserID:     user.UserID,
		LastSeen:   time.Now().UnixMilli(),
		DeviceName: deviceName,
	}
	device, err := s.Queries.CreateDevice(ctx, deviceParams)
	if err != nil {
		zap.L().Error("could not create device", zap.Error(err))
		return "", ErrInternal
	}

	// create JWT
	_, tokenString, _ := s.TokenAuth.Encode(map[string]interface{}{"user_id": user.UserID, "device_id": device.DeviceID})

	return tokenString, nil
}

func (s *UserService) LogoutUser(ctx context.Context) error {
	_, _, deviceID, err := extractFromJWT(ctx)
	if err != nil {
		zap.L().Error("could not logout user")
		return err
	}
	deviceUUID := &uuid.UUID{}
	deviceUUID.Scan(deviceID)

	err = s.Queries.DeleteDevice(ctx, *deviceUUID)
	if err != nil {
		zap.L().Error("could not delete existing device", zap.Error(err))
		return ErrInternal
	}
	return nil
}

func (s *UserService) CheckJWT(ctx context.Context) (string, string, error) {
	_, userID, deviceID, err := extractFromJWT(ctx)
	if err != nil {
		return "", "", err
	}
	deviceUUID := &uuid.UUID{}
	err = deviceUUID.Scan(deviceID)
	if err != nil {
		zap.L().Error("could not convert deviceID of type string to pgtype.UUID", zap.Error(err))
		return "", "", ErrInternal
	}
	userUUID := &uuid.UUID{}
	err = userUUID.Scan(userID)
	if err != nil {
		zap.L().Error("could not convert userID of type string to pgtype.UUID", zap.Error(err))
		return "", "", ErrInternal
	}

	query := dbqueries.GetDeviceByIDAndUserIDParams{
		DeviceID: *deviceUUID,
		UserID:   *userUUID,
	}

	device, err := s.Queries.GetDeviceByIDAndUserID(ctx, query)
	if err != nil {
		zap.L().Warn("Device with this deviceID for user with this userID not found", zap.String("DeviceID", deviceID), zap.String("UserID", userID))
		return "", "", ErrDeviceNotFound
	}
	if device.UserID != *userUUID {
		zap.L().Error("userID of jwt did not match userID of the provided device", zap.String("userID-jwt", userID), zap.String("userID-device", device.UserID.String()))
		return "", "", ErrInternal
	}

	return userID, deviceID, nil
}

func extractFromJWT(ctx context.Context) (jwt.Token, string, string, error) {
	token, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		zap.L().Error("JWT not valid or other error from verifier middleware", zap.Error(err))
		return nil, "", "", err
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		zap.L().Error("claims did not contain user_id or user_id was no string")
		return nil, "", "", ErrInternal
	}
	deviceID, ok := claims["device_id"].(string)
	if !ok {
		zap.L().Error("claims did not contain device_id or device_id was no string")
		return nil, "", "", ErrInternal
	}
	return token, userID, deviceID, nil
}
