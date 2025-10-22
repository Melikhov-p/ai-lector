// Package app содержит слой оркестрации приложения (Application Layer)
package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/user"
)

// UserApp - Application Service для работы с пользователями
// Координирует взаимодействие между доменными сервисами и внешними слоями
type UserApp struct {
	userService *user.Service
}

// NewUserApp создаёт новый экземпляр UserApp
func NewUserApp(userService *user.Service) *UserApp {
	return &UserApp{
		userService: userService,
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

// SearchUser универсальный метод поиска пользователя по различным критериям
// Поддерживает поиск по телефону и email (расширяемо для других критериев)
func (a *UserApp) SearchUser(ctx context.Context, filters user.UserFilter) ([]*user.User, error) {
	const op = "app.UserApp.FindUser"

	var (
		usrs []*user.User
		err  error
	)

	usrs, err = a.userService.SearchUser(ctx, filters)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s failed to find user by %v with %w", op, filters, err)
	}

	return usrs, nil
}

func (a *UserApp) GetUserByID(ctx context.Context, userID int) (*user.User, error) {
	const op = "app.UserApp.GetUserByID"

	var (
		usr *user.User
		err error
	)

	usr, err = a.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by ID with %w", op, err)
	}

	return usr, nil
}

func (a *UserApp) GetUsersList(ctx context.Context) ([]*user.User, error) {
	const op = "app.UserApp.GetUsersList"

	var (
		usrs []*user.User
		err  error
	)

	usrs, err = a.userService.GetUsersList(ctx)

	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s failed to find user by %v with %w", op, err)
	}

	return usrs, nil
}

func (a *UserApp) Auth(ctx context.Context, phone, password string) (*user.User, error) {
	const op = "app.UserApp.Auth"

	var (
		usr *user.User
		err error
	)

	usr, err = a.userService.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by phone with %w", op, user.ErrInvalidPhone)
	}

	err = a.userService.ComparePassword(ctx, usr, password)
	if err != nil {
		return nil, fmt.Errorf("%s failed to compare user password with %w", op, err)
	}

	return usr, nil
}
