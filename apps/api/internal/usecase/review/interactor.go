package review

import (
	"context"
	"log/slog"

	"gogym-api/internal/adapter/dto"
)

type interactor struct {
	repo Repository
}

func NewInteractor(repo Repository) ReviewUseCase {
	return &interactor{
		repo: repo,
	}
}

// GetReviewsByGymID retrieves reviews for a specific gym with pagination
func (i *interactor) GetReviewsByGymID(ctx context.Context, gymID int64, cursor string, limit int) (*dto.GetReviewsResponse, error) {
	slog.InfoContext(ctx, "GetReviewsByGymID UseCase", "gym_id", gymID, "cursor", cursor, "limit", limit)

	// デフォルト値設定
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	// リポジトリからドメインエンティティを取得
	domainReviews, nextCursor, err := i.repo.GetByGymID(ctx, gymID, cursor, limit)
	if err != nil {
		return nil, err
	}

	// ドメインエンティティをDTOに変換してレスポンス作成
	response := dto.ToReviewsResponse(domainReviews, nextCursor)

	return response, nil
}
