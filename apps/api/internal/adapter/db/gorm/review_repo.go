// internal/adapter/db/gorm/gym_repo.go
// 役割: ジムドメイン用リポジトリ実装（Infrastructure Layer）
// record↔entity変換を行いDB操作する実装。ドメインエンティティとGORMレコード間の変換を担当
package gorm

import (
	"context"
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/review"
	reviewUsecase "gogym-api/internal/usecase/review"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) reviewUsecase.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) GetByGymID(ctx context.Context, id int64, cursor string, limit int) ([]review.Review, string, error) {
	var reviews []record.ReviewRecord
	query := r.db.WithContext(ctx).
		Model(&record.ReviewRecord{}).
		Where("gym_id = ?", id).
		Order("created_at DESC").
		Limit(limit)
	if cursor != "" {
		query = query.Where("created_at < ?", cursor)
	}
	if err := query.Find(&reviews).Error; err != nil {
		return nil, "", err
	}

	// レビューをドメインエンティティに変換
	var reviewEntities []review.Review
	for _, r := range reviews {
		reviewEntities = append(reviewEntities, *ToReviewEntity(&r))
	}
	var nextCursor string
	if len(reviews) > 0 {
		nextCursor = reviews[len(reviews)-1].CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return reviewEntities, nextCursor, nil
}
