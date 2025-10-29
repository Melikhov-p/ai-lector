package middlewares

import (
	"net/http"
	"strconv"

	"github.com/Melikhov-p/ai-lector/internal/consts"
	"github.com/go-chi/chi/v5"
)

func (mw *Middleware) IsOwner(handler http.Handler) http.Handler {
	return mw.WithAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(consts.UserIDContextKey)
		if userID == nil {
			mw.logger.Warn("userID missing in context")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		resourceOwnerIDStr := chi.URLParam(r, consts.UserIDURLParam.String())
		resourceOwnerID, err := strconv.Atoi(resourceOwnerIDStr)
		if err != nil {
			mw.logger.Warn("user id in URL param is missing")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userIDinteger, ok := userID.(int64)
		if !ok {
			mw.logger.Warn("user id in URL param is invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if userIDinteger != int64(resourceOwnerID) {
			mw.logger.Warn("user id from URL param not equal to token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		handler.ServeHTTP(w, r)
	}))
}
