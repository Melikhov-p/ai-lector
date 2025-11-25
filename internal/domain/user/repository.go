package user

import (
	"context"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/subscription"
)

// Repository определяет интерфейс для работы с хранилищем пользователей
type Repository interface {
	// SaveUser сохраняет нового пользователя
	SaveUser(ctx context.Context, user *User) (int64, error)

	// GetUserByID возвращает пользователя по внешнему ID
	GetUserByID(ctx context.Context, userID int64) (*User, error)

	GetUserByPhone(ctx context.Context, phone string) (*User, error)

	SearchUser(ctx context.Context, filter UserFilter) ([]*User, error)

	GetAllUsers(ctx context.Context) ([]*User, error)

	AddInterestToUser(ctx context.Context, userID int64, inter *interest.Interest) error

	GetUserInterestsByID(ctx context.Context, userID int64) ([]int64, error)

	UpdateUser(ctx context.Context, user *User) error
	SubscribeUser(ctx context.Context, usrSub *subscription.UserSubscription) error
}

// ListOptions опции для получения списка пользователей
type ListOptions struct {
	Limit  int    // Количество записей на странице
	Offset int    // Смещение
	SortBy string // Поле для сортировки (id, email, name, created_at)
	Order  string // Порядок сортировки (asc, desc)
}
