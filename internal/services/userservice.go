package services

import (
	"context"

	"github.com/alexedwards/argon2id"
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

func (s *UserService) CreateUser(ctx context.Context, username string, passwd string, role int32) (dbqueries.User, error) {
	hash, err := argon2id.CreateHash(passwd, argon2id.DefaultParams)
	if err != nil {
		return dbqueries.User{}, err
	}
	params := dbqueries.CreateUserParams{
		Username: username,
		Pwhash:   hash,
		Role:     role,
	}

	user, err := s.Queries.CreateUser(ctx, params)
	if err != nil {
		return dbqueries.User{}, err
	}

	return user, nil
}
