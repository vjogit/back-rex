package auth

import (
	"github.com/go-chi/chi/v5"
)

func RouteAuth(r chi.Router) {
	r.Post("/logging", logging)
	r.Get("/refresh", refreshToken)
	r.Get("/me", me)
}
