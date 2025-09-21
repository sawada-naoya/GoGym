package review

import (
	"context"
	"gogym-api/internal/domain/review"
)

func (u *UseCase) GetReviewsByGymID(ctx context.Context, gymID int64, cursor string, limit int) ([]review.Review, string, error) {
	return u.reviewRepo.GetByGymID(ctx, gymID, cursor, limit)
}
