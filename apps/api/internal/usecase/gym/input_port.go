package gym

import (
	"context"
	"gogym-api/internal/adapter/http/dto"
	"gogym-api/internal/domain/gym"
)

// handler → usecase
type UseCase interface {
	GetGym(ctx context.Context, id gym.ID) (*dto.GymResponse, error)
	RecommendGyms(ctx context.Context) ([]dto.GymResponse, error)
}
