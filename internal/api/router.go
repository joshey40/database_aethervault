package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/joshey40/database_aethervault/internal/config"
	"github.com/joshey40/database_aethervault/internal/handler"
)

func CreateRouter(config *config.Config, uh *handler.UserHandler, sh *handler.SyncHandler, mw *Middlewares, tokenAuth *jwtauth.JWTAuth) http.Handler {
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
		r.Put("/register", uh.CreateUser)
		r.Post("/login", uh.LoginUser)
		r.Route("/logout", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(mw.Authenticator(tokenAuth))

			r.Post("/", uh.LogoutUser)
		})

		r.Route("/status", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(mw.Authenticator(tokenAuth))

			r.Get("/", uh.Status)
		})
	})

	r.Get("/alive", uh.Status)

	r.Route("/sync", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(mw.Authenticator(tokenAuth))

		r.Get("{version}", sh.UserPull)
		r.Post("{version}", sh.UserPush)
	})

	return r
}
