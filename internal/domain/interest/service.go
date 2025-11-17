package interest

import (
	"context"
	"fmt"
)

// Service сервисный слой интересов.
type Service struct {
	repo Repository
}

// NewService новый сервис интересов.
func NewService(r Repository) *Service {
	return &Service{repo: r}
}

// CreateInterest создание нового интереса
func (s *Service) CreateInterest(ctx context.Context, name, emoji string) (*Interest, error) {
	const op = "domain.Interest.CreateInterest"

	var (
		inter      *Interest
		newInterID int64
		err        error
	)

	inter = NewInterest(name, emoji)

	newInterID, err = s.repo.SaveInterest(ctx, inter)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	inter.SetID(newInterID)

	return inter, nil
}

// GetByID получение интереса по ID
func (s *Service) GetByID(ctx context.Context, id int64) (*Interest, error) {
	const op = "domain.Interest.GetByID"

	var (
		inter *Interest
		err   error
	)

	inter, err = s.repo.GetInterestByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return inter, nil
}

// GetByName получение интереса по названию
func (s *Service) GetByName(ctx context.Context, name string) (*Interest, error) {
	const op = "domain.Interest.GetByName"

	var (
		inter *Interest
		err   error
	)

	inter, err = s.repo.GetInterestByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return inter, nil
}

func (s *Service) GetAll(ctx context.Context) ([]*Interest, error) {
	const op = "domain.Interest.GetAll"

	var (
		inters []*Interest
		err    error
	)

	inters, err = s.repo.GetAllInterests(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return inters, nil
}
