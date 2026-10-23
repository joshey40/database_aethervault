package main

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
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
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// load config
	config, err := config.LoadConf("config.yml")
	if err != nil {
		log.Fatalf("Loading config failed, error: %v", err)
	}

	// create and setup logger
	log := logger.CreateLogger(config)
	defer log.Sync()
	zap.ReplaceGlobals(log)
	zap.RedirectStdLog(log)

	zap.L().Info("Welcome to the AETHERVAULT server")
	zap.L().Info("Have fun building your decks and managing your collection")

	// load env file and vars
	err = godotenv.Load()
	if err != nil {
		zap.L().Error("Loading Env file failed", zap.Error(err))
	}

	dbpasswd := os.Getenv("DBPASSWD")
	jwtSecret := os.Getenv("JWT_SECRET")

	// generate jokenAuth from secret
	tokenAuth := jwtauth.New("HS256", []byte(jwtSecret), nil)

	// create connection to database via connection pool. This ensures concurrency
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s", config.DB.Host, config.DB.Port, config.DB.User, config.DB.Name, dbpasswd)
	pool, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		zap.L().Fatal("Creating db pool failed", zap.Error(err))
	}

	// create sql queries based upon the connection pool
	queries := dbqueries.New(pool)

	// create db layout
	goose.SetBaseFS(embeddedMigrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		zap.L().Fatal("Setting goose dialect failed", zap.Error(err))
	}
	db := stdlib.OpenDBFromPool(pool)
	err = goose.Up(db, "internal/db/migrations")
	if err != nil {
		zap.L().Fatal("Creating DB structure failed", zap.Error(err))
	}
	err = db.Close()
	if err != nil {
		zap.L().Fatal("Closing the initial db connection failed", zap.Error(err))
	}
	zap.L().Info("Initial DB structure created")

	// create services
	userService := services.NewUserService(pool, queries, tokenAuth)

	// create handlers
	userHandler := handler.NewUserHandler(userService)

	// create middlewares
	middlewares := api.NewMiddlewares(userService)

	// creating  the router
	r := api.CreateRouter(config, userHandler, middlewares, tokenAuth)
	listenString := fmt.Sprintf(":%d", config.API.Port)

	// creating the server
	server := http.Server{
		Addr:    listenString,
		Handler: r,
	}

	// starting the server
	go func() {
		zap.L().Info("Server started and listening", zap.Int("Port", config.API.Port))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("server stopped for weird reason", zap.Error(err))
		}
	}()

	<-ctx.Done()
	zap.L().Info("shutdown signal received")
	stop()

	exitCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	zap.L().Info("shutting down http server")
	err = server.Shutdown(exitCtx)
	if err != nil {
		log.Fatal("server could not be stopped gracefully", zap.Error(err))
	}
}
