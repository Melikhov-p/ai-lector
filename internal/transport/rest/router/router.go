package router

import (
	"fmt"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/consts"
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

	api.With(mw.WithoutAuth).Post("/login", handler.ForUser.Login) // Аутентификация
	api.With(mw.WithAuth).Post("/logout", handler.ForUser.Logout)  // Разлогиниться

	api.Route("/users", func(r chi.Router) { // корень пользовательских эндпоинтов
		r.Get("/", handler.ForUser.UserList)                         // получить всех пользователей
		r.With(mw.WithoutAuth).Post("/", handler.ForUser.CreateUser) // создать пользователя
		r.Get("/search", handler.ForUser.SearchUser)                 // найти пользователя

		r.With(mw.IsOwner).Route(buildPatternWithID(consts.UserIDURLParam), func(r chi.Router) { // корень эндпоинтов конкретного пользователя
			r.Get("/", handler.ForUser.GetByID)       // получить конкретного пользователя
			r.Route("/interest", func(r chi.Router) { // корень эндпоинтов для интересов конкретного пользователя
				r.Post(buildPatternWithID(consts.InterestIDURLParam), handler.ForUser.AddInterests) // добавить интерес пользователю
			})
		})
	})

	api.Route("/interests", func(r chi.Router) {
		r.Get("/", handler.ForInterest.GetInterestsList)
		r.Get(buildPatternWithID(consts.InterestIDURLParam), handler.ForInterest.GetInterest)
	})

	r.Mount("/api", api)
	return r
}

// buildPatternWithID возвращает концовку эндпоинта с consts.URLParam для какого либо ID
//
/*
 Example:
 In: consts.UserIDURLParam
 Return: /{userID}
*/
func buildPatternWithID(param consts.URLParam) string {
	return fmt.Sprintf("/{%s}", param)
}
