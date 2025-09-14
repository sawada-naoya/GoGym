package review

import (
	"context"
	"gogym-api/internal/domain/review"
)

type UseCase struct {
	reviewRepo ReviewRepository
}

func NewUseCase(reviewRepo ReviewRepository) *UseCase {
	return &UseCase{
		reviewRepo: reviewRepo,
	}
}

type ReviewRepository interface {
	GetByGymID(ctx context.Context, gymID int64, cursor string, limit int) ([]review.Review, string, error)
}
