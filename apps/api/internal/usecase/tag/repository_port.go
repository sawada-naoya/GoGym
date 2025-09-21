package tag

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// usecase → repository
type Repository interface {
	FindAll(ctx context.Context) ([]gym.Tag, error)
	FindByIDs(ctx context.Context, ids []gym.ID) ([]gym.Tag, error)
	FindByNames(ctx context.Context, names []string) ([]gym.Tag, error)
	Create(ctx context.Context, tag *gym.Tag) error
	CreateMany(ctx context.Context, tags []gym.Tag) error
}