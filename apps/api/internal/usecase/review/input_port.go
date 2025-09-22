package review

import (
	"context"
	"gogym-api/internal/adapter/http/dto"
)

// handler → usecase
type UseCase interface {
	GetReviewsByGymID(ctx context.Context, gymID int64, cursor string, limit int) (*dto.GetReviewsResponse, error)
}