package mapper

import (
	"gogym-api/internal/adapter/repository/record"
	"gogym-api/internal/domain/gym"
	"gogym-api/internal/domain/review"
)

func ToReviewEntity(r *record.ReviewRecord) *review.Review {
	rating, _ := review.NewRating(r.Rating)

	var comment *string
	if r.Comment != nil && *r.Comment != "" {
		comment = r.Comment
	}

	entity := &review.Review{
		BaseEntity: gym.BaseEntity{
			ID:        gym.ID(r.ID), // レビューIDを使用、ジムIDではない
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
		UserID:          gym.ID(r.UserID),
		GymID:           gym.ID(r.GymID),
		Rating:          rating,
		Comment:         comment,
		UserDisplayName: nil, // レコードに存在しない
	}

	return entity
}

// FromReviewEntity はReviewドメインエンティティをReviewRecordに変換する
func FromReviewEntity(r *review.Review) *record.ReviewRecord {

	record := &record.ReviewRecord{
		ID:      int64(r.ID),
		UserID:  int64(r.UserID),
		GymID:   int64(r.GymID),
		Rating:  r.Rating.Int(),
		Comment: r.Comment,
	}

	return record
}
