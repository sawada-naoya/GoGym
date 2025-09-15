// internal/adapter/db/gorm/gym_mapper.go
// 役割: Review Entity ↔ Record 変換ユーティリティ（Infrastructure Layer）
// ジムドメインエンティティとGORMレコード間の双方向変換を担当
package gorm

import (
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/gym"
	"gogym-api/internal/domain/review"
)

func ToReviewEntity(r *record.ReviewRecord) *review.Review {
	rating, _ := review.NewRating(r.Rating)

	var comment *string
	if r.Content != "" {
		comment = &r.Content
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
	// コメントの処理（nullable）
	var content string
	if r.Comment != nil {
		content = *r.Comment
	}

	record := &record.ReviewRecord{
		ID:      int64(r.ID),
		Title:   "", // エンティティに存在しない - コメントから導出可能
		Content: content,
		Rating:  r.Rating.Int(),
		GymID:   int64(r.GymID),
		UserID:  int64(r.UserID),
	}

	return record
}
