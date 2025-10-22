package dto

// CreateUserDTO - DTO для создания пользователя
type CreateUserDTO struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}

// AuthUserDTO авторизация пользователя
type AuthUserDTO struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// UserDTO - DTO пользователя
type UserDTO struct {
	ID        int64  `json:"id"`
	Phone     string `json:"phone"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}

// UsersDTO - DTO массива пользователей
type UsersDTO struct {
	Users []*UserDTO `json:"users"`
}

// SearchUserDTO - DTO поиска пользователя
type SearchUserDTO struct {
	Phone     string `json:"phone"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
}
