package memory

import (
	"context"

	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/google/uuid"
)

type Storage struct {
	users         map[uuid.UUID]*user.User
	nextUserIndex int64
}

func NewStorage() *Storage {
	return &Storage{
		users:         make(map[uuid.UUID]*user.User),
		nextUserIndex: 1,
	}
}

func (s *Storage) Create(ctx context.Context, user *user.User) error {
	user.SetID(s.nextUserIndex)
	s.users[user.UUID()] = user

	s.nextUserIndex++
	return nil
}

func (s *Storage) GetByUUID(ctx context.Context, uuid uuid.UUID) (*user.User, error) {
	if usr, ok := s.users[uuid]; ok {
		return usr, nil
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	for _, v := range s.users {
		if v.Phone() == phone {
			return v, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) Search(ctx context.Context, filter user.UserFilter) ([]*user.User, error) {
	match := make([]*user.User, 0)
	finded := make(map[uuid.UUID]struct{})

	for _, v := range s.users {
		if _, ok := finded[v.UUID()]; !ok && filter.Email != "" && filter.Email == v.Email() {
			match = append(match, v)
			finded[v.UUID()] = struct{}{}
		}
	}
	for _, v := range s.users {
		if _, ok := finded[v.UUID()]; !ok && filter.Phone != "" && filter.Phone == v.Phone() {
			match = append(match, v)
			finded[v.UUID()] = struct{}{}
		}
	}
	for _, v := range s.users {
		if _, ok := finded[v.UUID()]; !ok && filter.FirstName != "" && filter.FirstName == v.FirstName() {
			match = append(match, v)
			finded[v.UUID()] = struct{}{}
		}
	}
	for _, v := range s.users {
		if _, ok := finded[v.UUID()]; !ok && filter.LastName != "" && filter.LastName == v.LastName() {
			match = append(match, v)
			finded[v.UUID()] = struct{}{}
		}
	}

	return match, nil
}

func (s *Storage) GetByID(ctx context.Context, userID int) (*user.User, error) {
	for _, v := range s.users {
		if v.ID() == int64(userID) {
			return v, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) GetAllUsers(ctx context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0, len(s.users))

	for _, v := range s.users {
		res = append(res, v)
	}

	return res, nil
}
