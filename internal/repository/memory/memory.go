package memory

import (
	"context"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
	"github.com/Melikhov-p/ai-lector/internal/domain/user"
	"github.com/google/uuid"
)

type Storage struct {
	users             map[uuid.UUID]*user.User
	interests         map[int64]*interest.Interest
	userInterests     map[int64][]*interest.Interest
	nextUserIndex     int64
	nextInterestIndex int64
}

func NewStorage() *Storage {
	str := &Storage{
		users:             make(map[uuid.UUID]*user.User),
		interests:         make(map[int64]*interest.Interest),
		userInterests:     make(map[int64][]*interest.Interest),
		nextUserIndex:     1,
		nextInterestIndex: 1,
	}

	return str
}

func (s *Storage) SaveUser(_ context.Context, user *user.User) error {
	user.SetID(s.nextUserIndex)
	s.users[user.UUID()] = user

	s.nextUserIndex++
	return nil
}

func (s *Storage) UpdateUser(_ context.Context, user *user.User) error {
	s.users[user.UUID()] = user

	return nil
}

func (s *Storage) GetUserByUUID(_ context.Context, uuid uuid.UUID) (*user.User, error) {
	if usr, ok := s.users[uuid]; ok {
		return usr, nil
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) GetUserByPhone(_ context.Context, phone string) (*user.User, error) {
	for _, v := range s.users {
		if v.Phone() == phone {
			return v, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) SearchUser(_ context.Context, filter user.UserFilter) ([]*user.User, error) {
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

func (s *Storage) GetUserByID(_ context.Context, userID int64) (*user.User, error) {
	for _, v := range s.users {
		if v.ID() == userID {
			return v, nil
		}
	}

	return nil, user.ErrUserNotFound
}

func (s *Storage) GetAllUsers(_ context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0, len(s.users))

	for _, v := range s.users {
		res = append(res, v)
	}

	return res, nil
}

func (s *Storage) SaveInterest(_ context.Context, interest *interest.Interest) (int64, error) {
	interest.SetID(s.nextInterestIndex)
	s.interests[s.nextInterestIndex] = interest
	s.nextInterestIndex++

	return s.nextInterestIndex, nil
}

func (s *Storage) DeleteInterest(_ context.Context, id int64) error {
	delete(s.interests, id)

	return nil
}

func (s *Storage) GetInterestByID(_ context.Context, id int64) (*interest.Interest, error) {
	if inter, ok := s.interests[id]; ok {
		return inter, nil
	}

	return nil, interest.ErrInterestNotFound
}

func (s *Storage) GetInterestByName(_ context.Context, name string) (*interest.Interest, error) {
	for _, v := range s.interests {
		if v.Name() == name {
			return v, nil
		}
	}

	return nil, interest.ErrInterestNotFound
}

func (s *Storage) AddInterestToUser(_ context.Context, userID int64, inter *interest.Interest) error {
	s.userInterests[userID] = append(s.userInterests[userID], inter)

	return nil
}

func (s *Storage) GetAllInterests(_ context.Context) ([]*interest.Interest, error) {
	res := make([]*interest.Interest, 0)

	for _, v := range s.interests {
		res = append(res, v)
	}

	return res, nil
}

func (s *Storage) GetUserInterestsByID(_ context.Context, userID int64) ([]int64, error) {
	var res []int64

	if inter, ok := s.userInterests[userID]; ok {
		for _, v := range inter {
			res = append(res, v.ID())
		}

		return res, nil
	}

	return nil, user.ErrUserInterestsEmpty
}
