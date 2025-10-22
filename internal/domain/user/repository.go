package user

import (
	"context"
)

// Repository определяет интерфейс для работы с хранилищем пользователей
type Repository interface {
	// Create создаёт нового пользователя
	Create(ctx context.Context, user *User) error

	// GetByID возвращает пользователя по внешнему ID
	GetByID(ctx context.Context, userID int) (*User, error)

	// // GetByEmail возвращает пользователя по email
	// GetByEmail(ctx context.Context, email string) (*User, error)

	GetByPhone(ctx context.Context, phone string) (*User, error)

	Search(ctx context.Context, filter UserFilter) ([]*User, error)

	GetAllUsers(ctx context.Context) ([]*User, error)

	// // Update обновляет данные пользователя
	// Update(ctx context.Context, user *User) error

	// // Delete удаляет пользователя по ID
	// Delete(ctx context.Context, id int64) error
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
