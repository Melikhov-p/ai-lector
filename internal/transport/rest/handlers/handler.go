package handlers

import (
	"github.com/Melikhov-p/ai-lector/internal/app"
	"go.uber.org/zap"
)

type Handlers struct {
	log         *zap.Logger
	ForUser     *userHandlers
	ForInterest *interestHandlers
}

func NewHandlers(l *zap.Logger, a *app.UserApp, i *app.InterestApp) *Handlers {
	return &Handlers{
		log:         l,
		ForUser:     newUserHandlers(l, a, i),
		ForInterest: newInterestHandlers(l, i),
	}
}
