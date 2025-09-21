package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// handler → usecase
type UseCase interface {
	GetGym(ctx context.Context, id gym.ID) (*gym.Gym, error)
	RecommendGyms(ctx context.Context) ([]gym.Gym, error)
}
