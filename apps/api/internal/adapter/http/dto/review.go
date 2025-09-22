package dto

import (
	dom "gogym-api/internal/domain/review"
)

type Review struct {
	ID        int64  `json:"id"`
	GymID     int64  `json:"gym_id"`
	UserID    int64  `json:"user_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type GetReviewsResponse struct {
	Reviews []Review `json:"reviews"`
	Cursor  string   `json:"cursor"`
}

// ToReviewDTO converts domain Review to DTO Review
func ToReviewDTO(domainReview dom.Review) Review {
	comment := ""
	if domainReview.Comment != nil {
		comment = *domainReview.Comment
	}

	return Review{
		ID:        int64(domainReview.ID),
		GymID:     int64(domainReview.GymID),
		UserID:    int64(domainReview.UserID),
		Rating:    domainReview.Rating.Int(),
		Comment:   comment,
		CreatedAt: domainReview.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// ToReviewsResponse converts domain Reviews and cursor to response DTO
func ToReviewsResponse(domainReviews []dom.Review, nextCursor string) *GetReviewsResponse {
	reviews := make([]Review, len(domainReviews))
	for i, domainReview := range domainReviews {
		reviews[i] = ToReviewDTO(domainReview)
	}

	return &GetReviewsResponse{
		Reviews: reviews,
		Cursor:  nextCursor,
	}
}