package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joshey40/database_aethervault/internal/config"
	"github.com/joshey40/database_aethervault/internal/handler"
)

func CreateRouter(config *config.Config, h *handler.UserHandler, tokenAuth *jwtauth.JWTAuth) http.Handler {
	// creating new Router
	r := chi.NewRouter()

	// registering base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)

	// adding "/auth" route to router to handle user login, logout and register
	r.Route("/auth", func(r chi.Router) {
		r.Put("/register", h.CreateUser)
		r.Post("/login", h.LoginUser)
		r.Route("/logout", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator(tokenAuth))

			r.Post("/", h.LogoutUser)

		})
	})

	return r
}
