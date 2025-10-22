package main

import (
	"net/http"

	app "github.com/Melikhov-p/ai-lector/internal/app/user"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/repository/memory"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/handlers"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/router"
	"github.com/Melikhov-p/ai-lector/pkg/logger"
)

func main() {
	log, err := logger.BuildLogger("DEBUG")
	if err != nil {
		panic("Logger")
	}

	usrService := user.NewService(memory.NewStorage())

	userApp := app.NewUserApp(usrService)

	mux := router.NewRouter(handlers.NewHandlers(log, userApp), log)

	log.Debug("Run server")

	if err = http.ListenAndServe(":8080", mux); err != nil {
		panic("server panic")
	}
}
