package main

import (
	"gitlab.sudovi.me/erp/file-ms-api/handlers"
	"gitlab.sudovi.me/erp/file-ms-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/oykos-development-hub/celeritas"
)

func routes(app *celeritas.Celeritas, middleware *middleware.Middleware, handlers *handlers.Handlers) *chi.Mux {
	// middleware must come before any routes

	// Kreirajte novi router
	r := chi.NewRouter()

	// Konfigurišite CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(cors.Handler)

	r.Route("/api", func(rt chi.Router) {
		rt.Post("/files", handlers.FileHandler.CreateFile)
		rt.Get("/files/{id}", handlers.FileHandler.GetFileById)
		rt.Delete("/files/{id}", handlers.FileHandler.DeleteFile)
		rt.Get("/download/*", handlers.FileHandler.GetFile)
		rt.Get("/file-overview/{id}", handlers.FileHandler.FileOverview)

		rt.Get("/read-articles", handlers.FileHandler.ReadArticles)
	})

	return r
}
