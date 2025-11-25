// Package app содержит слой оркестрации приложения (Application Layer)
package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/subscription"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
	"go.uber.org/zap"
)

// UserApp - Application Service для работы с пользователями
// Координирует взаимодействие между доменными сервисами и внешними слоями
type UserApp struct {
	userService *user.Service
	interestApp *InterestApp
	subApp      *SubscriptionApp
	log         *zap.Logger
}

// NewUserApp создаёт новый экземпляр UserApp
func NewUserApp(userService *user.Service, a *InterestApp, l *zap.Logger) *UserApp {
	return &UserApp{
		userService: userService,
		interestApp: a,
		log:         l,
	}
}

// CreateUser создаёт нового пользователя (use case)
func (a *UserApp) CreateUser(ctx context.Context, phone, firstName, password string) (*user.User, error) {
	const op = "app.UserApp"

	var (
		usr *user.User
		err error
	)

	// Используем новый универсальный метод поиска
	usr, err = a.userService.GetUserByPhone(ctx, phone)
	if err != nil && !errors.Is(err, user.ErrUserNotFound) {
		return nil, fmt.Errorf("%s failed to check user existence with %w", op, err)
	}

	if usr != nil {
		return nil, user.ErrPhoneAlreadyExist
	}

	usr, err = a.userService.CreateUser(ctx, phone, firstName, password)
	if err != nil {
		return nil, fmt.Errorf("%s failed to create user with %w", op, err)
	}

	return usr, err
}

// SearchUser универсальный метод поиска пользователя по различным критериям.
// Поддерживает поиск по телефону и email (расширяемо для других критериев)
func (a *UserApp) SearchUser(ctx context.Context, filters user.UserFilter) ([]*user.User, error) {
	const op = "app.UserApp.FindUser"

	var (
		usrs      []*user.User
		intersIDs []int64
		inters    []*interest.Interest
		err       error
	)

	usrs, err = a.userService.SearchUser(ctx, filters)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s failed to find user by %v with %w", op, filters, err)
	}

	for _, u := range usrs {
		intersIDs, err = a.userService.GetInterests(ctx, u)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			a.log.Error("Failed to get user interests", zap.Int64("UserID", u.ID()), zap.Error(err))
		}

		inters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
		if err != nil {
			a.log.Error("Failed to get interests", zap.Int64("UserID", u.ID()), zap.Error(err))
		}

		u.SetInterests(inters)
	}

	return usrs, nil
}

func (a *UserApp) GetUserByID(ctx context.Context, userID int64) (*user.User, error) {
	const op = "app.UserApp.GetUserByID"

	var (
		usr       *user.User
		intersIDs []int64
		inters    []*interest.Interest
		err       error
	)

	usr, err = a.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by ID with %w", op, err)
	}

	intersIDs, err = a.userService.GetInterests(ctx, usr)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		a.log.Error("Failed to get user interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	inters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
	if err != nil {
		a.log.Error("Failed to get interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	usr.SetInterests(inters)

	return usr, nil
}

func (a *UserApp) GetUsersList(ctx context.Context) ([]*user.User, error) {
	const op = "app.UserApp.GetUsersList"

	var (
		usrs      []*user.User
		intersIDs []int64
		inters    []*interest.Interest
		err       error
	)

	usrs, err = a.userService.GetUsersList(ctx)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s failed to find user with %w", op, err)
	}

	for _, u := range usrs {
		intersIDs, err = a.userService.GetInterests(ctx, u)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			a.log.Error("Failed to get user interests", zap.Int64("UserID", u.ID()), zap.Error(err))
		}

		inters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
		if err != nil {
			a.log.Error("Failed to get interests", zap.Int64("UserID", u.ID()), zap.Error(err))
		}

		u.SetInterests(inters)
	}

	return usrs, nil
}

func (a *UserApp) Auth(ctx context.Context, phone, password string) (*user.User, error) {
	const op = "app.UserApp.Auth"

	var (
		usr       *user.User
		intersIDs []int64
		inters    []*interest.Interest
		err       error
	)

	usr, err = a.userService.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by phone with %w", op, user.ErrInvalidPhone)
	}

	err = a.userService.ComparePassword(usr, password)
	if err != nil {
		return nil, fmt.Errorf("%s failed to compare user password with %w", op, err)
	}

	intersIDs, err = a.userService.GetInterests(ctx, usr)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		a.log.Error("Failed to get user interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	inters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
	if err != nil {
		a.log.Error("Failed to get interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	usr.SetInterests(inters)

	return usr, nil
}

func (a *UserApp) AddInterest(ctx context.Context, userID, interestID int64) error {
	const op = "app.UserApp.AddInterest"

	var (
		inter      *interest.Interest
		intersIDs  []int64
		userInters []*interest.Interest
		usr        *user.User
		err        error
	)

	usr, err = a.userService.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s failed to get user by ID %w", op, err)
	}

	intersIDs, err = a.userService.GetInterests(ctx, usr)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		a.log.Error("Failed to get user interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	userInters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
	if err != nil {
		a.log.Error("Failed to get interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	usr.SetInterests(userInters)

	if userHasInterest := a.userService.CheckUserInterestByID(usr, interestID); userHasInterest {
		return user.ErrInterestAlreadyExists
	}

	inter, err = a.interestApp.GetInterestByID(ctx, interestID)
	if err != nil {
		return fmt.Errorf("%s failed to get interest by id %w", op, err)
	}

	err = a.userService.AddInterest(ctx, usr, inter)
	if err != nil {
		return fmt.Errorf("%s failed to add interest to user %w", op, err)
	}

	return nil
}

func (a *UserApp) UpdateUser(ctx context.Context, userID int64, inDTO *dto.UpdateUserDTO) (*user.User, error) {
	const op = "app.UserApp.UpdateUser"

	var (
		usr       *user.User
		intersIDs []int64
		inters    []*interest.Interest
		err       error
	)

	usr, err = a.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by ID %w", op, err)
	}

	err = a.userService.UpdateUser(ctx, usr, inDTO)
	if err != nil {
		return nil, fmt.Errorf("%s failed to update user %w", op, err)
	}

	intersIDs, err = a.userService.GetInterests(ctx, usr)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		a.log.Error("Failed to get user interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	inters, err = a.interestApp.GetInterestsBatchByIDs(ctx, intersIDs)
	if err != nil {
		a.log.Error("Failed to get interests", zap.Int64("UserID", usr.ID()), zap.Error(err))
	}

	usr.SetInterests(inters)

	return usr, nil
}

// DenyTrialRequest убрать 1 пробный запрос у пользователя
func (a *UserApp) DenyTrialRequest(ctx context.Context, usr *user.User) error {
	const op = "app.UserApp.removeTrialRequest"

	err := a.userService.DenyTrialRequest(ctx, usr, 1)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *UserApp) Subscribe(ctx context.Context, userID, subID int64) error {
	const op = "app.UserApp.subscribe"

	var (
		usr *user.User
		sub *subscription.Subscription
		err error
	)

	usr, err = a.userService.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s failed to get user by ID %w", op, err)
	}

	sub, err = a.subApp.GetSubscriptionByID(ctx, subID)
	if err != nil {
		return fmt.Errorf("%s failed to get subscription by ID %w", op, err)
	}

	err = a.userService.Subscribe(ctx, usr, sub)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
