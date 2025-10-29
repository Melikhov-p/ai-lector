package main

import (
	"context"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/app"
	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/repository/memory"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/handlers"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/router"
	"github.com/Melikhov-p/ai-lector/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	log, err := logger.BuildLogger("DEBUG")
	if err != nil {
		panic("Logger")
	}

	store := memory.NewStorage()
	createInterests(store)

	interestService := interest.NewService(store)
	interestApp := app.NewInterestApp(interestService)

	usrService := user.NewService(store)
	userApp := app.NewUserApp(usrService, interestApp, log)
	_, err = userApp.CreateUser(context.Background(), "89997776655", "Павел", "1q2w3e")
	if err != nil {
		log.Error("error creating user", zap.Error(err))
	}

	mux := router.NewRouter(handlers.NewHandlers(log, userApp, interestApp), log)

	log.Debug("Run server")

	if err = http.ListenAndServe(":8080", mux); err != nil {
		panic("server panic")
	}
}

func createInterests(store *memory.Storage) {
	_ = store.SaveInterest(context.Background(), interest.NewInterest("футбол", "⚽️"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("компьютерные игры", "🖥"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("гимнастика", "🤸🏻‍♂️"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("кино", "🎥"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("скейтборд", "🛹"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("музыка", "🎶"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("машины", "🚘"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("хоккей", "🏒"))
	_ = store.SaveInterest(context.Background(), interest.NewInterest("баскетбол", "🏀"))
}
