package app

import (
	"context"
	"fmt"

	"github.com/Melikhov-p/ai-lector/internal/domain/subscription"
)

type SubscriptionApp struct {
	service *subscription.Service
}

func NewSubscriptionApp(s *subscription.Service) *SubscriptionApp {
	return &SubscriptionApp{
		service: s,
	}
}

func (sa *SubscriptionApp) GetSubscriptionByID(ctx context.Context, id int64) (*subscription.Subscription, error) {
	const op = "app.subscription.GetSubscriptionByID"

	var (
		sub *subscription.Subscription
		err error
	)

	sub, err = sa.service.GetSubscriptionByID(context.Background(), id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}

func (sa *SubscriptionApp) GetSubscriptionByName(name string) (*subscription.Subscription, error) {
	const op = "app.subscription.getSubscriptionByName"

	var (
		sub *subscription.Subscription
		err error
	)

	sub, err = sa.service.GetSubscriptionByName(context.Background(), name)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return sub, nil
}
