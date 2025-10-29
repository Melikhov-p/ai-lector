package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/transport/rest/dto"
	"github.com/Melikhov-p/ai-lector/internal/util"
	"golang.org/x/crypto/bcrypt"
)

// Service предоставляет доменные сервисы для работы с пользователями
// Domain Service инкапсулирует бизнес-логику создания, обновления и управления пользователями
type Service struct {
	repo Repository
}

// NewService создаёт новый экземпляр Service
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateUser создаёт нового пользователя с валидацией бизнес-правил
func (s *Service) CreateUser(ctx context.Context, phone, firstName, password string) (*User, error) {
	const op = "domain.Service.CreateUser"

	var (
		normPhone string
		usr       *User
		passHash  []byte
		err       error
	)

	if firstName == "" {
		return nil, ErrEmptyName
	}

	normPhone, err = util.NormalizePhone(phone)
	if err != nil {
		return nil, fmt.Errorf("%s failed to normalize phone with %w", op, err)
	}

	passHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s failed to generate hash for password with %w", op, err)
	}

	// Создание пользователя через domain model
	usr, err = NewUser(normPhone, firstName, passHash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Сохранение в репозиторий
	if err = s.repo.SaveUser(ctx, usr); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return usr, nil
}

// GetUserByPhone возвращает пользователя по номеру телефона или error ErrUserNotFound
func (s *Service) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	const op = "domain.Service.GetUserByPhone"

	var (
		usr       *User
		normPhone string
		err       error
	)

	normPhone, err = util.NormalizePhone(phone)
	if err != nil {
		return nil, fmt.Errorf("%s failed to normalize phone with %w", op, err)
	}

	usr, err = s.repo.GetUserByPhone(ctx, normPhone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s failed to get user by phone with %w", op, err)
	}

	return usr, nil
}

func (s *Service) SearchUser(ctx context.Context, filter UserFilter) ([]*User, error) {
	const op = "domain.Service.SearchUser"

	var (
		usrs []*User
		err  error
	)

	usrs, err = s.repo.SearchUser(ctx, filter)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf("%s failed to search users with %w", op, err)
	}

	return usrs, nil
}

func (s *Service) GetByID(ctx context.Context, userID int64) (*User, error) {
	const op = "domain.Service.User.GetByID"

	var (
		usr *User
		err error
	)

	usr, err = s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get user by id with %w", op, err)
	}

	return usr, nil
}

func (s *Service) GetUsersList(ctx context.Context) ([]*User, error) {
	const op = "domain.Service.User.GetUsersList"

	var (
		usrs []*User
		err  error
	)

	usrs, err = s.repo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get list of users with %w", op, err)
	}

	return usrs, nil
}

func (s *Service) ComparePassword(usr *User, password string) error {
	err := bcrypt.CompareHashAndPassword(usr.passHash, []byte(password))

	return err
}

func (s *Service) AddInterest(ctx context.Context, usr *User, inter *interest.Interest) error {
	const op = "domain.User.Service.AddInterest"

	err := s.repo.AddInterestToUser(ctx, usr.ID(), inter)
	if err != nil {
		return fmt.Errorf("%s failed to add interest to user %w", op, err)
	}

	return nil
}

func (s *Service) GetInterests(ctx context.Context, usr *User) ([]*interest.Interest, error) {
	const op = "domain.UserService.GetInterests"

	var (
		inters []*interest.Interest
		err    error
	)

	inters, err = s.repo.GetUserInterestsByID(ctx, usr.ID())
	if err != nil {
		return nil, fmt.Errorf("%s failed to get interests for user %w", op, err)
	}

	return inters, nil
}

func (s *Service) UpdateUser(ctx context.Context, usr *User, inDTO *dto.UpdateUserDTO) error {
	const op = "domain.UserService.UpdateUser"

	var (
		passHash []byte
		err      error
	)

	if inDTO.Phone != "" {
		usr.phone = inDTO.Phone
	}
	if inDTO.Email != "" {
		usr.email = inDTO.Email
	}
	if inDTO.FirstName != "" {
		usr.firstName = inDTO.FirstName
	}
	if inDTO.LastName != "" {
		usr.lastName = inDTO.LastName
	}
	if inDTO.Password != "" {
		passHash, err = bcrypt.GenerateFromPassword([]byte(inDTO.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("%s failed to generate hash for password with %w", op, err)
		}

		usr.passHash = passHash
	}

	err = s.repo.UpdateUser(ctx, usr)
	if err != nil {
		return fmt.Errorf("%s failed to save user with %w", op, err)
	}

	return nil
}
