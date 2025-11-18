package dto

// CreateUserDTO - DTO для создания пользователя
type CreateUserDTO struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	Password  string `json:"password"`
}

// UpdateUserDTO - DTO для обновления пользователя
type UpdateUserDTO struct {
	Phone     string `json:"phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Class     int    `json:"class"`
}

// AuthUserDTO авторизация пользователя
type AuthUserDTO struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

// InterestDTO интерес
type InterestDTO struct {
	Id    int64  `json:"id"`    // ID - идентификатор
	Name  string `json:"name"`  // Name - название интереса
	Emoji string `json:"emoji"` // Emoji - смайлик интереса
}

type InterestListDTO struct {
	Interests []*InterestDTO `json:"interests"`
}

// UserDTO - DTO пользователя
type UserDTO struct {
	ID            int64          `json:"id"`
	Phone         string         `json:"phone"`
	Email         string         `json:"email,omitempty"`
	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name,omitempty"`
	TrialRequests int            `json:"trial_requests,omitempty"`
	Interests     []*InterestDTO `json:"interests,omitempty"`
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

type ExplainDTO struct {
	Explanation string `json:"explanation"`
}
