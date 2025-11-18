package handlers

import (
	"github.com/Melikhov-p/ai-lector/internal/app"
	"go.uber.org/zap"
)

type Handlers struct {
	log         *zap.Logger
	ForUser     *userHandlers
	ForInterest *interestHandlers
	ForExplain  *explainHandlers
}

func NewHandlers(l *zap.Logger, a *app.UserApp, i *app.InterestApp, c aiClient) *Handlers {
	return &Handlers{
		log:         l,
		ForUser:     newUserHandlers(l, a, i),
		ForInterest: newInterestHandlers(l, i),
		ForExplain:  newExplaneHandlers(l, a, c),
	}
}
