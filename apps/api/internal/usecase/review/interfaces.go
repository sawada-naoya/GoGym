package review

import (
	"context"
)

// ReviewUseCase defines the interface for review business logic
type ReviewUseCase interface {
	// TODO: Create review
	CreateReview(ctx context.Context, req CreateReviewRequest) (*ReviewResponse, error)
	// TODO: Get review by ID
	GetReview(ctx context.Context, reviewID int64) (*ReviewResponse, error)
	// TODO: Get reviews for gym
	GetReviewsByGym(ctx context.Context, gymID int64, req GetReviewsRequest) (*GetReviewsResponse, error)
	// TODO: Update review
	UpdateReview(ctx context.Context, reviewID int64, req UpdateReviewRequest) (*ReviewResponse, error)
	// TODO: Delete review
	DeleteReview(ctx context.Context, reviewID int64) error
}

// Repository interfaces
type ReviewRepository interface {
	// TODO: Repository methods
}

// Request/Response DTOs
type CreateReviewRequest struct {
	// TODO: Add fields
}

type UpdateReviewRequest struct {
	// TODO: Add fields
}

type GetReviewsRequest struct {
	// TODO: Add pagination, filtering
}

type ReviewResponse struct {
	// TODO: Add response fields
}

type GetReviewsResponse struct {
	// TODO: Add reviews list and pagination
}