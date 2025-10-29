package handlers

import (
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
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type userHandlers struct {
	log         *zap.Logger
	userApp     *app.UserApp
	interestApp *app.InterestApp
}

func newUserHandlers(l *zap.Logger, a *app.UserApp, i *app.InterestApp) *userHandlers {
	return &userHandlers{
		log:         l,
		userApp:     a,
		interestApp: i,
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err = uh.userApp.CreateUser(r.Context(), inDTO.Phone, inDTO.FirstName, inDTO.Password)
	if err != nil {
		uh.log.Error("error while creating new user", zap.Any("CreateDTO", inDTO), zap.Error(err))

		if errors.Is(err, user.ErrPhoneAlreadyExist) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(outDTO); err != nil {
		uh.log.Error("error while encoding outDTO", zap.Any("outDTO", outDTO), zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err = strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err = uh.userApp.UpdateUser(r.Context(), userID, &inDTO)
	if err != nil {
		uh.log.Error("error while updating user", zap.Any("UpdateUser", inDTO), zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	if err = json.NewEncoder(w).Encode(outDTO); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err = uh.userApp.GetUserByID(r.Context(), int64(id))
	if err != nil {
		uh.log.Warn("searching for user are not success", zap.Error(err))
		if errors.Is(err, user.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	if err = json.NewEncoder(w).Encode(&outDTO); err != nil {
		uh.log.Error("error while encoding to user DTO", zap.Error(err), zap.Any("DTO", outDTO))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, usr := range usrs {
		usrDTO := mapper.FromUserToDTO(usr)
		outDTO.Users = append(outDTO.Users, &usrDTO)
	}

	if err = json.NewEncoder(w).Encode(outDTO); err != nil {
		uh.log.Error("error while encoding out DTO", zap.Error(err), zap.Any("DTO", outDTO))
		w.WriteHeader(http.StatusInternalServerError)
	}
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
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, usr := range usrs {
		usrDTO := mapper.FromUserToDTO(usr)
		outDTO.Users = append(outDTO.Users, &usrDTO)
	}

	if err = json.NewEncoder(w).Encode(outDTO); err != nil {
		uh.log.Error("error while decoding out DTO", zap.Error(err), zap.Any("DTO", outDTO))
		w.WriteHeader(http.StatusInternalServerError)
	}
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err = uh.userApp.Auth(r.Context(), inDTO.Phone, inDTO.Password)
	if err != nil {
		uh.log.Warn("failed to authenticate user", zap.Error(err), zap.Any("inDTO", inDTO))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	outDTO = mapper.FromUserToDTO(usr)
	token, err = auth.BuildJWTToken(usr, "supersecretkey", 24*time.Hour)
	if err != nil {
		uh.log.Error("error while building jwt token", zap.Error(err), zap.Any("User", usr))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	})

	if err = json.NewEncoder(w).Encode(&outDTO); err != nil {
		uh.log.Error("error while encoding out DTO", zap.Error(err), zap.Any("DTO", outDTO))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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

	w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	interestIDParam = chi.URLParam(r, consts.InterestIDURLParam.String())
	interestID, err = strconv.Atoi(interestIDParam)
	if err != nil {
		uh.log.Error("invalid interestID", zap.Int(consts.InterestIDURLParam.String(), interestID))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.userApp.AddInterest(r.Context(), userID, int64(interestID))
	if err != nil {
		uh.log.Warn("failed to add interest", zap.Error(err))

		if errors.Is(err, user.ErrInterestAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		if errors.Is(err, user.ErrUserNotFound) || errors.Is(err, user.ErrInterestNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
