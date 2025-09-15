package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

func (uc *UseCase) GetGym(ctx context.Context, id gym.ID) (*gym.Gym, error) {
	uc.logger.InfoContext(ctx, "getting gym for search/preview", "gym_id", id)

	if id == 0 {
		return nil, gym.NewDomainError(gym.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	foundGym, err := uc.gymRepo.FindByID(ctx, id)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get gym", "gym_id", id, "error", err)
		return nil, gym.NewDomainErrorWithCause(err, "gym_not_found", "gym not found")
	}

	reviewStats, err := uc.gymRepo.GetReviewStats(ctx, id)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get review stats", "gym_id", id, "error", err)
	} else {
		foundGym.AverageRating = reviewStats.AverageRating
		foundGym.ReviewCount = reviewStats.ReviewCount
	}

	return foundGym, nil
}
