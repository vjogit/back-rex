// This example demonstrates how to serve static files from your filesystem.
//
// Boot the server:
//
//	$ go run main.go
//
// Client requests:
//
//	$ curl http://localhost:3333/files/
package main

import (
	"back-rex/pkg/auth"
	"back-rex/pkg/feedback"
	"back-rex/pkg/utils"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

const ConnString = "host=localhost port=5432 user=postgres dbname=rex-ema  password=root sslmode=disable"

type SecutityCtx struct{}

func Security(role string) func(next http.Handler) http.Handler {

	security := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			user, err := auth.GetUserFromCookie(r)

			if err != nil {
				render.Render(w, r, utils.ErrUnauthorizedRequest(err))
				return
			}

			if !user.Roles.Valid || !strings.Contains(user.Roles.String, role) {
				render.Render(w, r, utils.ErrUnauthorizedRequest(errors.New("bad permissions")))
				return
			}

			requestCtx := context.WithValue(r.Context(), SecutityCtx{}, user)
			r = r.WithContext(requestCtx)

			next.ServeHTTP(w, r)
		})
	}
	return security
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Route("/api/v0", func(r chi.Router) {
		r.Use(utils.MakeDatabaseMiddleware(ConnString))

		r.Route("/auth", auth.RouteAuth)
		r.Route("/feedback", feedback.RouteFeedback)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})

	http.ListenAndServe(":3333", r)
}
