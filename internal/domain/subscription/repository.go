package subscription

import "context"

type repository interface {
	CreateSubscription(ctx context.Context, subscription *Subscription) error
	GetSubscriptionByName(ctx context.Context, name string) (*Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int64) (*Subscription, error)
}
