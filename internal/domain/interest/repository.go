package interest

import "context"

type Repository interface {
	SaveInterest(ctx context.Context, interest *Interest) (int64, error)
	DeleteInterest(ctx context.Context, id int64) error
	GetInterestByID(ctx context.Context, id int64) (*Interest, error)
	GetInterestByName(ctx context.Context, name string) (*Interest, error)
	GetAllInterests(ctx context.Context) ([]*Interest, error)
}
