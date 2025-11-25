package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Melikhov-p/ai-lector/internal/app"
	"github.com/Melikhov-p/ai-lector/internal/auth"
	"github.com/Melikhov-p/ai-lector/internal/consts"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/mapper"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/responser"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type userHandlers struct {
	log         *zap.Logger
	userApp     *app.UserApp
	interestApp *app.InterestApp
	resp        *responser.Responser
}

func newUserHandlers(l *zap.Logger, a *app.UserApp, i *app.InterestApp) *userHandlers {
	return &userHandlers{
		log:         l,
		userApp:     a,
		interestApp: i,
		resp:        responser.NewResponser(l),
	}
}

// CreateUser хэндлер создания пользователя
func (uh *userHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var (
		usr    *user.User
		inDTO  dto.CreateUserDTO
		outDTO dto.UserDTO
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&inDTO); err != nil {
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	usr, err = uh.userApp.CreateUser(r.Context(), inDTO.Phone, inDTO.FirstName, inDTO.Password)
	if err != nil {
		uh.log.Error("error while creating new user", zap.Any("CreateDTO", inDTO), zap.Error(err))

		if errors.Is(err, user.ErrPhoneAlreadyExist) {
			uh.resp.WriteError(w, http.StatusConflict, user.ErrPhoneAlreadyExist)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)

	uh.resp.WriteJSON(w, http.StatusCreated, outDTO)
}

func (uh *userHandlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var (
		usr    *user.User
		inDTO  dto.UpdateUserDTO
		outDTO dto.UserDTO
		userID int64
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&inDTO); err != nil {
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	userID, err = strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	usr, err = uh.userApp.UpdateUser(r.Context(), userID, &inDTO)
	if err != nil {
		uh.log.Error("error while updating user", zap.Any("UpdateUser", inDTO), zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			uh.resp.WriteError(w, http.StatusNotFound, user.ErrUserNotFound)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	uh.resp.WriteJSON(w, http.StatusOK, outDTO)
}

func (uh *userHandlers) GetByID(w http.ResponseWriter, r *http.Request) {
	var (
		usr    *user.User
		idStr  string
		id     int
		outDTO dto.UserDTO
		err    error
	)

	idStr = chi.URLParam(r, consts.UserIDURLParam.String())
	if id, err = strconv.Atoi(idStr); err != nil {
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	usr, err = uh.userApp.GetUserByID(r.Context(), int64(id))
	if err != nil {
		uh.log.Warn("searching for user are not success", zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			uh.resp.WriteError(w, http.StatusNotFound, user.ErrUserNotFound)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	uh.resp.WriteJSON(w, http.StatusOK, outDTO)
}

// UserList получение списка пользователей.
func (uh *userHandlers) UserList(w http.ResponseWriter, r *http.Request) {
	var (
		usrs   []*user.User
		outDTO dto.UsersDTO
		err    error
	)

	usrs, err = uh.userApp.GetUsersList(r.Context())
	if err != nil {
		uh.log.Warn("getting users list are not success", zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			uh.resp.WriteError(w, http.StatusNotFound, user.ErrUserNotFound)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	for _, usr := range usrs {
		usrDTO := mapper.FromUserToDTO(usr)
		outDTO.Users = append(outDTO.Users, &usrDTO)
	}

	uh.resp.WriteJSON(w, http.StatusOK, outDTO)
}

// SearchUser поиск пользователя
func (uh *userHandlers) SearchUser(w http.ResponseWriter, r *http.Request) {
	var (
		usrs      []*user.User
		outDTO    dto.UsersDTO
		usrFilter user.UserFilter
		err       error
	)

	usrFilter = user.UserFilter{
		Email:     r.URL.Query().Get("email"),
		Phone:     r.URL.Query().Get("phone"),
		FirstName: r.URL.Query().Get("firstName"),
		LastName:  r.URL.Query().Get("lastName"),
	}

	usrs, err = uh.userApp.SearchUser(r.Context(), usrFilter)
	if err != nil {
		uh.log.Warn("searching for user are not success", zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			uh.resp.WriteError(w, http.StatusNotFound, user.ErrUserNotFound)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	for _, usr := range usrs {
		usrDTO := mapper.FromUserToDTO(usr)
		outDTO.Users = append(outDTO.Users, &usrDTO)
	}

	uh.resp.WriteJSON(w, http.StatusOK, outDTO)
}

// Login аутентификация
func (uh *userHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var (
		inDTO  dto.AuthUserDTO
		outDTO dto.UserDTO
		usr    *user.User
		token  string
		err    error
	)

	if err = json.NewDecoder(r.Body).Decode(&inDTO); err != nil {
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	usr, err = uh.userApp.Auth(r.Context(), inDTO.Phone, inDTO.Password)
	if err != nil {
		uh.log.Warn("failed to authenticate user", zap.Error(err), zap.Any("inDTO", inDTO))
		uh.resp.WriteError(w, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	token, err = auth.BuildJWTToken(usr, "supersecretkey", 24*time.Hour)
	if err != nil {
		uh.log.Error("error while building jwt token", zap.Error(err), zap.Any("User", usr))
		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	uh.resp.WriteJSON(w, http.StatusOK, outDTO)
}

// Logout выход
func (uh *userHandlers) Logout(w http.ResponseWriter, _ *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // удалить
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	uh.resp.WriteJSON(w, http.StatusOK, nil)
}

// AddInterests добавить пользователю интерес
func (uh *userHandlers) AddInterests(w http.ResponseWriter, r *http.Request) {
	var (
		interestIDParam string
		userID          int64
		interestID      int
		err             error
		ok              bool
	)

	userID, ok = r.Context().Value(consts.UserIDContextKey).(int64)
	if !ok {
		uh.log.Error("invalid user ID in context", zap.Int64("UserID", userID))
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	interestIDParam = chi.URLParam(r, consts.InterestIDURLParam.String())
	interestID, err = strconv.Atoi(interestIDParam)
	if err != nil {
		uh.log.Error("invalid interestID", zap.Int(consts.InterestIDURLParam.String(), interestID))
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	err = uh.userApp.AddInterest(r.Context(), userID, int64(interestID))
	if err != nil {
		uh.log.Warn("failed to add interest", zap.Error(err))

		if errors.Is(err, user.ErrInterestAlreadyExists) {
			uh.resp.WriteError(w, http.StatusConflict, user.ErrInterestAlreadyExists)
			return
		}

		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrInterestNotFound) {
			uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
			return
		}

		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	uh.resp.WriteJSON(w, http.StatusOK, nil)
}

func (uh *userHandlers) RemoveInterest(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("NOT IMPLEMENTED"))
}

// Subscribe оформить подписку
func (uh *userHandlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		userID              int64
		subscriptionIDParam string
		subscriptionID      int
		err                 error
		ok                  bool
	)

	userID, ok = r.Context().Value(consts.UserIDContextKey).(int64)
	if !ok {
		uh.log.Error("invalid user ID in context", zap.Int64("UserID", userID))
		uh.resp.WriteError(w, http.StatusBadRequest, ErrBadRequest)
		return
	}

	subscriptionIDParam = chi.URLParam(r, consts.SubscriptionIDParam.String())
	subscriptionID, err = strconv.Atoi(subscriptionIDParam)

	err = uh.userApp.Subscribe(context.Background(), userID, int64(subscriptionID))
	if err != nil {
		uh.log.Warn(
			"failed to subscribe to subscription",
			zap.Error(err),
			zap.Int64("UserID", userID),
			zap.Int("SubscriptionID", subscriptionID))
		uh.resp.WriteError(w, http.StatusInternalServerError, ErrInternalServerError)
		return
	}

	uh.resp.WriteJSON(w, http.StatusOK, nil)
}
