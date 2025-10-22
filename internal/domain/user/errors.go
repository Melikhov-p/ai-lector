package user

import "errors"

var (
	// ErrUserNotFound возвращается, когда пользователь не найден
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidEmail возвращается при невалидном формате email
	ErrInvalidEmail = errors.New("invalid email format")

	ErrInvalidPhone = errors.New("invalid phone")

	// ErrEmptyName возвращается, когда имя пустое
	ErrEmptyName = errors.New("name cannot be empty")

	// ErrEmailAlreadyExists возвращается, когда email уже существует
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInterestAlreadyExists возвращается, когда интерес уже существует
	ErrInterestAlreadyExists = errors.New("interest already exists")

	// ErrInterestNotFound возвращается, когда интерес не найден
	ErrInterestNotFound = errors.New("interest not found")

	// ErrPhonealreadyExist возвращается, когда пользователь с таким телефоном уже есть.
	ErrPhoneAlreadyExist = errors.New("user with this phone already exists")
	ErrInvalidPassword   = errors.New("invalid password")
)
