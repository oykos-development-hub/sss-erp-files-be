package main

import (
	"gitlab.sudovi.me/erp/file-ms-api/handlers"
	"gitlab.sudovi.me/erp/file-ms-api/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/oykos-development-hub/celeritas"
)

func routes(app *celeritas.Celeritas, middleware *middleware.Middleware, handlers *handlers.Handlers) *chi.Mux {
	// middleware must come before any routes

	//api
	app.Routes.Route("/api", func(rt chi.Router) {

		rt.Get("/file/*", handlers.FileHandler.ShowFile)

		rt.Post("/files", handlers.FileHandler.CreateFile)
		rt.Get("/files/{id}", handlers.FileHandler.GetFileById)
		rt.Get("/files", handlers.FileHandler.GetFileList)
		rt.Put("/files/{id}", handlers.FileHandler.UpdateFile)
		rt.Delete("/files/{id}", handlers.FileHandler.DeleteFile)
	})

	return app.Routes
}
