package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joshey40/database_aethervault/internal/config"
	"github.com/joshey40/database_aethervault/internal/handler"
	"github.com/joshey40/database_aethervault/internal/logger"
	"go.uber.org/zap"
)

func StartRouter(config *config.Config, h *handler.UserHandler) {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Put("/register/", h.CreateUser)
	})

	listenString := fmt.Sprintf(":%d", config.API.Port)
	logger.L().Info("Server started and listening", zap.Int("Port", config.API.Port))
	http.ListenAndServe(listenString, r)
}
