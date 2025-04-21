package router

import (
	"github.com/decipherhub/sui-graph/server/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	middleware "github.com/go-chi/chi/v5/middleware"
)

func SetupRouter() http.Handler {
	r := chi.NewRouter()

	// Global Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	// Health check
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Authenticated API routes
	r.Route("/api", func(r chi.Router) {
		//r.Use(authmiddle.APIKeyAuth)
		r.Get("/graph", handler.GetGraphByCheckpoint)
	})

	return r
}
