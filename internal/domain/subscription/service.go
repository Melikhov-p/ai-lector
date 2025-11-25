package subscription

import (
	"context"
	"fmt"
	"time"
)

// Service сервис подписок
type Service struct {
	repo repository
}

// NewService новый сервис подписок
func NewService(r repository) *Service {
	return &Service{repo: r}
}

// CreateSubscription создание новой подписки
func (s *Service) CreateSubscription(
	ctx context.Context,
	name string,
	price float64,
	duration time.Duration,
) (*Subscription, error) {
	const op = "subscription.service.CreateSubscription"

	var (
		sub *Subscription
		err error
	)

	sub = NewSubscription(name, price, duration)

	err = s.repo.CreateSubscription(ctx, sub)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

// GetSubscriptionByName получить подписку по имени
func (s *Service) GetSubscriptionByName(ctx context.Context, name string) (*Subscription, error) {
	const op = "subscription.service.GetSubscriptionByName"

	var (
		sub *Subscription
		err error
	)

	sub, err = s.repo.GetSubscriptionByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (s *Service) GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error) {
	const op = "subscription.service.GetSubscriptionByID"

	var (
		sub *Subscription
		err error
	)

	sub, err = s.repo.GetSubscriptionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}
