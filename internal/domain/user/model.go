// Package user содержит доменные модели и бизнес-логику пользователя
package user

import (
	"time"

	"github.com/google/uuid"
)

// User представляет пользователя системы
type User struct {
	id        int64     // Внутренний ID для БД
	uuid      uuid.UUID // Внешний UUID для API
	email     string
	phone     string
	passHash  []byte
	firstName string
	lastName  string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser создаёт нового пользователя (без ID, будет присвоен при сохранении)
func NewUser(phone, firstName string, passHash []byte) (*User, error) {
	return &User{
		uuid:      uuid.New(),
		phone:     phone,
		firstName: firstName,
		passHash:  passHash,
		createdAt: time.Now(),
		updatedAt: time.Now(),
	}, nil
}

// ID возвращает внутренний идентификатор пользователя
func (u *User) ID() int64 {
	return u.id
}

// UUID возвращает внешний UUID пользователя
func (u *User) UUID() uuid.UUID {
	return u.uuid
}

// Email возвращает email пользователя
func (u *User) Email() string {
	return u.email
}

// FirstName возвращает имя пользователя
func (u *User) FirstName() string {
	return u.firstName
}

// LastName возвращает имя пользователя
func (u *User) LastName() string {
	return u.lastName
}

// Phone возвращает номер телефона
func (u *User) Phone() string {
	return u.phone
}

// CreatedAt возвращает время создания пользователя
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt возвращает время последнего обновления пользователя
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// SetID устанавливает внутренний ID (используется репозиторием после сохранения)
func (u *User) SetID(id int64) {
	u.id = id
}

// UpdateName обновляет имя пользователя
func (u *User) UpdateFirstName(name string) error {
	if name == "" {
		return ErrEmptyName
	}
	u.firstName = name

	return nil
}

// UserFilter модель фильтра для поиска пользователя
type UserFilter struct {
	Email     string
	Phone     string
	FirstName string
	LastName  string
}
