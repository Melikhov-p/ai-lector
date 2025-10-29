package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/interest"
)

// InterestApp оркеструет сценарии интересов.
type InterestApp struct {
	interestService *interest.Service
}

// NewInterestApp новое приложение интересов.
func NewInterestApp(srv *interest.Service) *InterestApp {
	return &InterestApp{
		interestService: srv,
	}
}

// CreateInterest создать новый интерес
func (app *InterestApp) CreateInterest(ctx context.Context, name, emoji string) (*interest.Interest, error) {
	const op = "app.Interest.CreateInterest"

	var (
		inter *interest.Interest
		err   error
	)

	inter, err = app.interestService.GetByName(ctx, name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%s failed to check existing interest: %w", op, err)
	}

	if inter != nil {
		return nil, interest.ErrInterestAlreadyExist
	}

	inter, err = app.interestService.CreateInterest(ctx, name, emoji)
	if err != nil {
		return nil, fmt.Errorf("%s failed to create interest: %w", op, err)
	}

	return inter, nil
}

// GetInterestByID получить интерес по ID
func (app *InterestApp) GetInterestByID(ctx context.Context, interestID int64) (*interest.Interest, error) {
	const op = "app.Interest.GetInterestByID"

	var (
		inter *interest.Interest
		err   error
	)

	if inter, err = app.interestService.GetByID(ctx, interestID); err != nil {
		return nil, fmt.Errorf("%s failed to get interest by ID: %w", op, err)
	}

	return inter, nil
}

// GetInterests получить все интересы
func (app *InterestApp) GetInterests(ctx context.Context) ([]*interest.Interest, error) {
	const op = "app.Interest.GetInterests"

	var (
		interests []*interest.Interest
		err       error
	)

	interests, err = app.interestService.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s failed to get interests: %w", op, err)
	}

	return interests, nil
}
