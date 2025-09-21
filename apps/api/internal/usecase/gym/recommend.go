package gym

import (
	"context"
	"log/slog"
	"sort"

	"gogym-api/internal/domain/gym"
)

// RecommendGyms returns recommended gyms for top page
func (i *interactor) RecommendGyms(ctx context.Context) ([]gym.Gym, error) {
	slog.InfoContext(ctx, "RecommendGyms UseCase")

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

	return gyms, nil
}