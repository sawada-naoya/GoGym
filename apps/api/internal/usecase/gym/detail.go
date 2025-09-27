package gym

import (
	"context"

	"gogym-api/internal/adapter/http/dto"
	dom "gogym-api/internal/domain/gym"
)

// GetGym returns gym detail by ID
func (i *interactor) GetGym(ctx context.Context, id dom.ID) (*dto.GymResponse, error) {

	if id == 0 {
		return nil, dom.NewDomainError("invalid_gym_id")
	}

	foundGym, err := i.repo.FindByID(ctx, id)
	if err != nil {
		return nil, dom.NewDomainError("gym_not_found")
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
