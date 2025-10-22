package middlewares

import (
	"context"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/auth"
	"go.uber.org/zap"
)

func (mw *Middleware) WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userID      int64
			tokenCookie *http.Cookie
			token       string
			err         error
		)

		tokenCookie, err = r.Cookie("token")
		if err != nil {
			mw.logger.Warn("try to get access without token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token = tokenCookie.Value
		userID, err = auth.ValidateJWTToken(token, "supersecretkey")
		if err != nil {
			mw.logger.Warn("invalid token", zap.Error(err))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		h.ServeHTTP(w, r.WithContext(ctx))

	})
}
