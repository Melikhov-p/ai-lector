package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/subscription"
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
		newUserID int64
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
	if newUserID, err = s.repo.SaveUser(ctx, usr); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	usr.SetID(newUserID)

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

// GetInterests получить интересы пользователя
func (s *Service) GetInterests(ctx context.Context, usr *User) ([]int64, error) {
	const op = "user.service.GetInterests"

	var (
		inters []int64
		err    error
	)

	inters, err = s.repo.GetUserInterestsByID(ctx, usr.ID())
	if err != nil {
		return nil, fmt.Errorf("%s failed to get interests for user %w", op, err)
	}

	return inters, nil
}

// CheckUserInterestByID проверяет наличие у пользователя интереса с указанным ID
func (s *Service) CheckUserInterestByID(usr *User, interID int64) bool {
	for _, i := range usr.Interests() {
		if i.ID() == interID {
			return true
		}
	}

	return false
}

// UpdateUser обновить пользователя
func (s *Service) UpdateUser(ctx context.Context, usr *User, inDTO *dto.UpdateUserDTO) error {
	const op = "user.service.UpdateUser"

	var (
		passHash []byte
		err      error
	)

	if inDTO.Phone != "" {
		usr.phone = inDTO.Phone
	}
	if inDTO.Email != "" {
		if util.IsValidEmail(inDTO.Email) {
			usr.email = inDTO.Email
		} else {
			return ErrInvalidEmail
		}
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
	if inDTO.Age != 0 {
		usr.age = inDTO.Age
	}
	if inDTO.Class != 0 {
		usr.class = inDTO.Class
	}

	err = s.repo.UpdateUser(ctx, usr)
	if err != nil {
		return fmt.Errorf("%s failed to save user with %w", op, err)
	}

	return nil
}

// DenyTrialRequest убавить пробные запросы
func (s *Service) DenyTrialRequest(ctx context.Context, usr *User, count int) error {
	const op = "user.service.DenyTrialRequest"

	usr.trialRequests -= count

	err := s.repo.UpdateUser(ctx, usr)
	if err != nil {
		return fmt.Errorf("%s failed to update user with %w", op, err)
	}

	return nil
}

// Subscribe оформить подписку для пользователя
func (s *Service) Subscribe(ctx context.Context, usr *User, sub *subscription.Subscription) error {
	const op = "user.service.Subscribe"

	var (
		userSub *subscription.UserSubscription
		err     error
	)

	userSub = subscription.NewUserSubscription(usr.ID(), sub)

	err = s.repo.SubscribeUser(ctx, userSub)
	if err != nil {
		return fmt.Errorf("%s failed to subscribe to user with %w", op, err)
	}

	return nil
}
