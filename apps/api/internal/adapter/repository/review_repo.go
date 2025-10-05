package repository

import (
	"context"
	"gogym-api/internal/adapter/repository/mapper"
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/review"
	reviewUsecase "gogym-api/internal/usecase/review"

	"gorm.io/gorm"
)

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) reviewUsecase.Repository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) GetByGymID(ctx context.Context, id int64, cursor string, limit int) ([]review.Review, string, error) {
	var reviews []struct {
		record.ReviewRecord
		User record.User `gorm:"foreignKey:UserID"`
	}
	query := r.db.WithContext(ctx).
		Model(&record.ReviewRecord{}).
		Preload("User").
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
		entity := mapper.ToReviewEntity(&r.ReviewRecord)
		if r.User.Name != "" {
			entity.UserDisplayName = &r.User.Name
		}
		reviewEntities = append(reviewEntities, *entity)
	}
	var nextCursor string
	if len(reviews) > 0 {
		nextCursor = reviews[len(reviews)-1].CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}

	return reviewEntities, nextCursor, nil
}
