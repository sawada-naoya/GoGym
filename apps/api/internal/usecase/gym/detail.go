package gym

import (
	"context"
	"log/slog"

	"gogym-api/internal/adapter/http/dto"
	dom "gogym-api/internal/domain/gym"
)

// GetGym returns gym detail by ID
func (i *interactor) GetGym(ctx context.Context, id dom.ID) (*dto.GymResponse, error) {
	slog.InfoContext(ctx, "GetGym UseCase", "gym_id", id)

	if id == 0 {
		return nil, dom.NewDomainError(dom.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	foundGym, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return nil, dom.NewDomainErrorWithCause(err, "gym_not_found", "gym not found")
	}

	// レビュー統計を取得
	reviewStats, err := i.repo.GetReviewStatsForGyms(ctx, []dom.ID{id})
	if err == nil {
		if stats, exists := reviewStats[id]; exists {
			foundGym.AverageRating = stats.AverageRating
			foundGym.ReviewCount = stats.ReviewCount
		}
	}

	// DTOに変換して返却
	response := dto.ToGymResponse(*foundGym)
	return &response, nil
}
