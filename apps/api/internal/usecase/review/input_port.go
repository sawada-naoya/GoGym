package review

import (
	"context"
	"gogym-api/internal/adapter/dto"
)

// handler → usecase
type ReviewUseCase interface {
	GetReviewsByGymID(ctx context.Context, gymID int64, cursor string, limit int) (*dto.GetReviewsResponse, error)
}