package user

import (
	"context"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
)

// Repository определяет интерфейс для работы с хранилищем пользователей
type Repository interface {
	// SaveUser создаёт нового пользователя
	SaveUser(ctx context.Context, user *User) error

	// GetUserByID возвращает пользователя по внешнему ID
	GetUserByID(ctx context.Context, userID int64) (*User, error)

	GetUserByPhone(ctx context.Context, phone string) (*User, error)

	SearchUser(ctx context.Context, filter UserFilter) ([]*User, error)

	GetAllUsers(ctx context.Context) ([]*User, error)

	AddInterestToUser(ctx context.Context, userID int64, inter *interest.Interest) error

	GetUserInterestsByID(ctx context.Context, userID int64) ([]*interest.Interest, error)

	UpdateUser(ctx context.Context, user *User) error
}

// ListOptions опции для получения списка пользователей
type ListOptions struct {
	Limit  int    // Количество записей на странице
	Offset int    // Смещение
	SortBy string // Поле для сортировки (id, email, name, created_at)
	Order  string // Порядок сортировки (asc, desc)
}

// DefaultListOptions возвращает опции по умолчанию
func DefaultListOptions() ListOptions {
	return ListOptions{
		Limit:  20,
		Offset: 0,
		SortBy: "id",
		Order:  "asc",
	}
}
