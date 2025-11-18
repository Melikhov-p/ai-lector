package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/app"
	"github.com/Melikhov-p/ai-lector/internal/client/yandex"
	"github.com/Melikhov-p/ai-lector/internal/consts"
	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/repository/postgres"
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

	store, err := postgres.NewPostgresStorage("postgresql://ai_lector:password@localhost:5433/ailector_app")
	if err != nil {
		panic(err)
	}
	log.Debug("postgres running", zap.Any("store", store))

	//createInterests(store, log) НУЖНО ТОЛЬКО ПРИ ПЕРВОМ ЗАПУСКЕ
	//log.Debug("created interests")

	interestService := interest.NewService(store)
	interestApp := app.NewInterestApp(interestService)

	usrService := user.NewService(store)
	userApp := app.NewUserApp(usrService, interestApp, log)
	_, err = userApp.CreateUser(context.Background(), "89997776655", "Павел", "1q2w3e")
	if err != nil {
		log.Error("error creating user", zap.Error(err))
	}

	yandexClient := yandex.NewClient(
		"yandex",
		"",
		"https://llm.api.cloud.yandex.net/v1/chat/completions",
		"gpt://b1gm2onh41hvp1g7idqg/yandexgpt/latest",
	)
	log.Debug("yandex client running")

	yandexClient.SetInstruction(consts.Instruction)

	mux := router.NewRouter(handlers.NewHandlers(log, userApp, interestApp, yandexClient), log)

	log.Debug("Run server")

	if err = http.ListenAndServe(":8080", mux); err != nil {
		panic("server panic")
	}
}

func createInterests(store interest.Repository, log *zap.Logger) {
	ins := []*interest.Interest{
		interest.NewInterest("футбол", "⚽️"),
		interest.NewInterest("компьютерные игры", "🖥"),
		interest.NewInterest("гимнастика", "🤸🏻‍♂️"),
		interest.NewInterest("кино", "🎥"),
		interest.NewInterest("скейтборд", "🛹"),
		interest.NewInterest("музыка", "🎶"),
		interest.NewInterest("машины", "🚘"),
		interest.NewInterest("хоккей", "🏒"),
		interest.NewInterest("баскетбол", "🏀"),
	}

	for _, inter := range ins {
		newID, err := store.SaveInterest(context.Background(), inter)
		if err != nil {
			log.Error("", zap.Error(err))
		} else {
			fmt.Println("newInterestID:", newID)
		}
	}
}
