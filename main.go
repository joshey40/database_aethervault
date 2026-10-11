package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joshey40/database_aethervault/internal/api"
	"github.com/joshey40/database_aethervault/internal/config"
	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
	"github.com/joshey40/database_aethervault/internal/handler"
	"github.com/joshey40/database_aethervault/internal/logger"
	"github.com/joshey40/database_aethervault/internal/services"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

//go:embed internal/db/migrations/*.sql
var embeddedMigrations embed.FS

func main() {
	// create main context
	ctx := context.Background()

	// load config
	config, err := config.LoadConf("config.yml")
	if err != nil {
		logger.L().Error("Loading config failed", zap.Error(err))
	}

	// create connection to database via connection pool. This ensures concurrency
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=REPLACEME", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		logger.L().Fatal("Creating db pool failed", zap.Error(err))
	}

	// create sql queries based upon the connection pool
	queries := dbqueries.New(pool)

	// create db layout
	goose.SetBaseFS(embeddedMigrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		logger.L().Fatal("Setting goose dialect failed", zap.Error(err))
	}
	db := stdlib.OpenDBFromPool(pool)
	err = goose.Up(db, "internal/db/migrations")
	if err != nil {
		logger.L().Fatal("Creating DB structure failed", zap.Error(err))
	}
	err = db.Close()
	if err != nil {
		logger.L().Fatal("Closing the initial db connection failed", zap.Error(err))
	}
	logger.L().Info("Initial DB structure created")

	// create services
	userService := services.NewUserService(pool, queries)

	// create handlers
	userHandler := handler.NewUserHandler(userService)

	// creating and starting the router
	api.StartRouter(config, userHandler)
}
