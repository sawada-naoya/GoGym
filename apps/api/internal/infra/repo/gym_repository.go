package repo

import (
	"context"
	"gogym-api/internal/domain/gym"
)

type gymRepository struct{}

func NewGymRepository() *gymRepository {
	return &gymRepository{}
}

// FindByID は単一のジムを返す（ダミー）
func (r *gymRepository) FindByID(ctx context.Context, id gym.ID) (*gym.Gym, error) {
	return &gym.Gym{
		ID:          id,
		Name:        "ゴールドジム 渋谷東京",
		Description: "ダミーデータ: 本格的なトレーニング環境を提供するジム",
	}, nil
}

// Search はダミーのジム一覧を返す
func (r *gymRepository) Search(ctx context.Context, query gym.SearchQuery) (*gym.PaginatedResult[gym.Gym], error) {
	items := []gym.Gym{
		{ID: 1, Name: "ゴールドジム 渋谷東京", Description: "渋谷の中心にある人気ジム"},
		{ID: 2, Name: "ゴールドジム 原宿ANNEX", Description: "原宿の人気ジム"},
		{ID: 3, Name: "ゴールドジム 南青山", Description: "おしゃれな青山のジム"},
	}

	return &gym.PaginatedResult[gym.Gym]{Items: items, NextCursor: ""}, nil
}

// GetReviewStatsForGyms は全ジムを★4.5で返す
func (r *gymRepository) GetReviewStatsForGyms(ctx context.Context, gymIDs []gym.ID) (map[gym.ID]*gym.ReviewStats, error) {
	stats := make(map[gym.ID]*gym.ReviewStats)
	for _, id := range gymIDs {
		stats[id] = &gym.ReviewStats{
			AverageRating: func() *float32 { f := float32(4.5); return &f }(),
			ReviewCount:   10,
		}
	}
	return stats, nil
}