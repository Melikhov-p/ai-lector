package router

import (
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/middlewares"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/handlers"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func NewRouter(handler *handlers.Handlers, l *zap.Logger) *chi.Mux {
	r := chi.NewRouter()

	mw := middlewares.NewMiddleware(l)

	r.Use(mw.CorsMiddleware)
	r.Use(mw.WithLogging)

	api := chi.NewRouter()

	// Регистрируем маршруты
	api.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	api.Post("/auth", handler.ForUser.Auth)

	api.Route("/users", func(r chi.Router) {
		r.Get("/", handler.ForUser.UserList)
		r.Post("/", handler.ForUser.CreateUser)
		r.Get("/search", handler.ForUser.SearchUser)
		r.Get("/{id}", handler.ForUser.GetByID)
	})

	r.Mount("/api", api)
	return r
}
