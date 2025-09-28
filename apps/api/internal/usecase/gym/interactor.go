package gym

import (
	"log/slog"

	"gogym-api/internal/usecase/tag"

	"context"
	"sort"

	"gogym-api/internal/adapter/dto"
	"gogym-api/internal/domain/gym"
	dom "gogym-api/internal/domain/gym"
)

type interactor struct {
	repo    Repository
	tagRepo tag.Repository
	logger  *slog.Logger
}

func NewInteractor(repo Repository, tagRepo tag.Repository, logger *slog.Logger) GymUseCase {
	return &interactor{
		repo:    repo,
		tagRepo: tagRepo,
		logger:  logger,
	}
}

// RecommendGyms returns recommended gyms for top page
func (i *interactor) RecommendGyms(ctx context.Context) ([]dto.GymResponse, error) {

	// トップページ用の固定パラメータ
	searchQuery := gym.SearchQuery{
		Query:    "",
		Location: nil,
		RadiusM:  nil,
		Pagination: gym.Pagination{
			Cursor: "",
			Limit:  10, // トップページ用に10件固定
		},
	}

	result, err := i.repo.Search(ctx, searchQuery)
	if err != nil {
		return nil, err
	}

	gyms := result.Items

	// レビュー統計を取得してソート
	gymIDs := make([]gym.ID, len(gyms))
	for j, g := range gyms {
		gymIDs[j] = g.ID
	}

	reviewStats, err := i.repo.GetReviewStatsForGyms(ctx, gymIDs)
	if err == nil {
		for k := range gyms {
			if stats, exists := reviewStats[gyms[k].ID]; exists {
				gyms[k].AverageRating = stats.AverageRating
				gyms[k].ReviewCount = stats.ReviewCount
			}
		}
	}

	// 評価順でソート
	sort.Slice(gyms, func(i, j int) bool {
		if gyms[i].AverageRating != nil && gyms[j].AverageRating != nil {
			if *gyms[i].AverageRating != *gyms[j].AverageRating {
				return *gyms[i].AverageRating > *gyms[j].AverageRating
			}
			return gyms[i].ReviewCount > gyms[j].ReviewCount
		}
		if gyms[i].AverageRating != nil && gyms[j].AverageRating == nil {
			return true
		}
		if gyms[i].AverageRating == nil && gyms[j].AverageRating != nil {
			return false
		}
		return gyms[i].ID < gyms[j].ID
	})

	// DTOに変換して返却
	responses := make([]dto.GymResponse, len(gyms))
	for i, gym := range gyms {
		responses[i] = dto.ToGymResponse(gym)
	}

	return responses, nil
}

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
