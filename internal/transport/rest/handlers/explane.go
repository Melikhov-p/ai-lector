package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/Melikhov-p/ai-lector/internal/app"
	"github.com/Melikhov-p/ai-lector/internal/consts"
	"github.com/Melikhov-p/ai-lector/internal/domain/explain"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/responser"
	"go.uber.org/zap"
)

type aiClient interface {
	MakeRequest(req string) (string, error)
}

type explainHandlers struct {
	log     *zap.Logger
	userApp *app.UserApp
	client  aiClient
	resp    *responser.Responser
}

func newExplaneHandlers(log *zap.Logger, userApp *app.UserApp, client aiClient) *explainHandlers {
	return &explainHandlers{
		log:     log,
		userApp: userApp,
		resp:    responser.NewResponser(log),
		client:  client,
	}
}

// GetExplain получить обхяснение заданной темы по предмету с учетом интересов пользователя в запросе
func (eh *explainHandlers) GetExplain(w http.ResponseWriter, r *http.Request) {
	var (
		userID         int64
		ok             bool
		theme, subject string
		u              *user.User
		err            error
		outDTO         dto.ExplainDTO
	)

	userID, ok = r.Context().Value(consts.UserIDContextKey).(int64)
	if !ok {
		eh.log.Error("failed to get user ID from context", zap.Any("Context.userID", r.Context().Value(consts.UserIDContextKey)))
		eh.resp.WriteError(w, http.StatusInternalServerError, errors.New("lost user"))
		return
	}
	theme = r.URL.Query().Get("theme")
	subject = r.URL.Query().Get("subject")

	if theme == "" || subject == "" {
		eh.log.Error(
			"empty query params in request",
			zap.Any("user", userID),
			zap.Any("theme", theme),
			zap.Any("subject", subject),
		)
		eh.resp.WriteError(w, http.StatusBadRequest, errors.New("empty query params in request"))
		return
	}
	u, err = eh.userApp.GetUserByID(r.Context(), userID)
	if err != nil {
		eh.log.Error("failed to get user by id", zap.Int64("userID", userID))
		eh.resp.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	outDTO.Explanation, err = eh.client.MakeRequest(explain.GetExplainRequestString(u, theme, subject))
	if err != nil {
		eh.log.Error("failed to get explain", zap.Error(err), zap.Any("user", u))
		eh.resp.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if u.TrialRequests() > 0 { // убираем у пользователя один пробный запрос если он был
		go func(log *zap.Logger) {
			err = eh.userApp.DenyTrialRequest(context.Background(), u)
			if err != nil {
				log.Error("failed to deny trial request", zap.Error(err), zap.Any("user", u))
			}
		}(eh.log)
	}

	eh.resp.WriteJSON(w, http.StatusOK, &outDTO)
}
