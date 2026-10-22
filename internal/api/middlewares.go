package api

import (
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joshey40/database_aethervault/internal/services"
)

type Middlewares struct {
	s services.Authservice
}

func NewMiddlewares(service services.Authservice) *Middlewares {
	return &Middlewares{s: service}
}

func (mw *Middlewares) Authenticator(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			_, _, err := mw.s.CheckJWT(r.Context())

			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(hfn)
	}
}
