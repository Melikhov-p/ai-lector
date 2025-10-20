package middlewares

import (
	"net/http"

	"go.uber.org/zap"
)

type (
	Middleware struct {
		logger *zap.Logger
	}

	responseData struct {
		status int
		size   int
	}

	// loggerResponseWriter модифицированный писатель для логирование запросов/ответов.
	loggerResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func NewMiddleware(logger *zap.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}
