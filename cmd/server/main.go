package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joshey40/database_aethervault/internal/api"
	"github.com/joshey40/database_aethervault/internal/config"
	dbqueries "github.com/joshey40/database_aethervault/internal/db/gen"
	"github.com/joshey40/database_aethervault/internal/handler"
	"github.com/joshey40/database_aethervault/internal/logger"
	"github.com/joshey40/database_aethervault/internal/services"
	"go.uber.org/zap"
)

func main() {
	// create main context
	ctx := context.Background()

	// load config
	config, err := config.LoadConf("../../config.yml")
	if err != nil {
		logger.L().Error("Loading config failed", zap.Error(err))
	}

	// create connection to database via connection pool. This ensures concurrency
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=REPLACEME", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		logger.L().Error("Creating db pool failed", zap.Error(err))
	}

	// create sql queries based upon the connection pool
	queries := dbqueries.New(pool)

	// create db layout
	err = queries.CreateDB(ctx)
	if err != nil {
		logger.L().Error("Creating initial DB strucure failed", zap.Error(err))
	}
	logger.L().Info("Initial DB structure created")
	// create services
	userService := services.NewUserService(pool, queries)

	// create handlers
	userHandler := handler.NewUserHandler(userService)

	// creating and starting the router
	api.StartRouter(config, userHandler)
}
