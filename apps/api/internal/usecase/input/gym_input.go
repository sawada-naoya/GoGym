package gym

import (
	"context"
	"gogym-api/internal/adapter/dto"
	"gogym-api/internal/domain/gym"
)

// handler → usecase
type GymUseCase interface {
	GetGym(ctx context.Context, id gym.ID) (*dto.GymResponse, error)
	RecommendGyms(ctx context.Context) ([]dto.GymResponse, error)
}
