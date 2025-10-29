// Package user содержит доменные модели и бизнес-логику пользователя
package user

import (
	"time"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/google/uuid"
)

// User представляет пользователя системы
type User struct {
	id            int64     // Внутренний ID для БД
	uuid          uuid.UUID // Внешний UUID для API
	trialRequests int       // Количество пробных запросов
	email         string
	phone         string
	passHash      []byte
	firstName     string
	lastName      string
	age           int
	createdAt     time.Time
	updatedAt     time.Time
	interests     []*interest.Interest
}

// NewUser создаёт нового пользователя (без ID, будет присвоен при сохранении)
func NewUser(phone, firstName string, passHash []byte) (*User, error) {
	return &User{
		uuid:          uuid.New(),
		phone:         phone,
		firstName:     firstName,
		passHash:      passHash,
		createdAt:     time.Now(),
		updatedAt:     time.Now(),
		trialRequests: 1,
		age:           0,
	}, nil
}

func NewUserFromDB(
	id int64,
	uuid uuid.UUID,
	trialRequests int,
	email string,
	phone string,
	passHash []byte,
	firstName string,
	lastName string,
	age int,
	createdAt time.Time,
	updatedAt time.Time,
) *User {
	return &User{
		id:            id,
		uuid:          uuid,
		trialRequests: trialRequests,
		email:         email,
		phone:         phone,
		passHash:      passHash,
		firstName:     firstName,
		lastName:      lastName,
		age:           age,
		createdAt:     createdAt,
		updatedAt:     updatedAt,
	}
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

func (u *User) PassHash() []byte {
	return u.passHash
}

// Age возраст пользователя
func (u *User) Age() int {
	return u.age
}

// SetAge установить возраст клиента
func (u *User) SetAge(age int) {
	u.age = age
}

// TrialRequests возвращает количество пробных запросов пользователя
func (u *User) TrialRequests() int {
	return u.trialRequests
}

// Interests возвращает интересы пользователя
func (u *User) Interests() []*interest.Interest {
	return u.interests
}

// SetInterests установить пользователю интересы
func (u *User) SetInterests(inters []*interest.Interest) {
	u.interests = inters
}

// SetID устанавливает внутренний ID (используется репозиторием после сохранения)
func (u *User) SetID(id int64) {
	u.id = id
}

// UpdateFirstName обновляет имя пользователя
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
