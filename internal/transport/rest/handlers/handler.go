package handlers

import (
	app "github.com/Melikhov-p/ai-lector/internal/app/user"
	"go.uber.org/zap"
)

type Handlers struct {
	log     *zap.Logger
	ForUser *UserHandlers
}

func NewHandlers(l *zap.Logger, a *app.UserApp) *Handlers {
	return &Handlers{
		log:     l,
		ForUser: newUserHandlers(l, a),
	}
}
