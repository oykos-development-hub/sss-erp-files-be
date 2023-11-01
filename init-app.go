package main

import (
	"log"
	"os"

	"gitlab.sudovi.me/erp/file-ms-api/handlers"
	"gitlab.sudovi.me/erp/file-ms-api/middleware"

	"github.com/oykos-development-hub/celeritas"
	"gitlab.sudovi.me/erp/file-ms-api/data"
	"gitlab.sudovi.me/erp/file-ms-api/services"
)

func initApplication() *celeritas.Celeritas {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// init celeritas
	cel := &celeritas.Celeritas{}
	err = cel.New(path)
	if err != nil {
		log.Fatal(err)
	}

	cel.AppName = "gitlab.sudovi.me/erp/file-ms-api"

	models := data.New(cel.DB.Pool)

	FileService := services.NewFileServiceImpl(cel, models.File)
	FileHandler := handlers.NewFileHandler(cel, FileService)

	myHandlers := &handlers.Handlers{FileHandler: FileHandler}

	myMiddleware := &middleware.Middleware{
		App: cel,
	}

	cel.Routes = routes(cel, myMiddleware, myHandlers)

	return cel
}
