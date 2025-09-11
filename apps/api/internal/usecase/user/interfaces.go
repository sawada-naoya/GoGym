package user

import (
	"context"
)

// UserUseCase defines the interface for user business logic
type UserUseCase interface {
	// TODO: User registration
	Register(ctx context.Context, req RegisterUserRequest) (*UserResponse, error)
	// TODO: User login
	Login(ctx context.Context, req LoginUserRequest) (*LoginUserResponse, error)
	// TODO: Get user profile
	GetUser(ctx context.Context, userID int64) (*UserResponse, error)
	// TODO: Update user profile
	UpdateUser(ctx context.Context, userID int64, req UpdateUserRequest) (*UserResponse, error)
	// TODO: Get user's favorite gyms
	GetUserFavorites(ctx context.Context, userID int64) (*UserFavoritesResponse, error)
	// TODO: Get user's reviews
	GetUserReviews(ctx context.Context, userID int64) (*UserReviewsResponse, error)
}

// Repository interface
type UserRepository interface {
	// TODO: Repository methods
}

// Request/Response DTOs
type RegisterUserRequest struct {
	// TODO: Add registration fields
}

type LoginUserRequest struct {
	// TODO: Add login fields
}

type UpdateUserRequest struct {
	// TODO: Add update fields
}

type UserResponse struct {
	// TODO: Add user response fields
}

type LoginUserResponse struct {
	// TODO: Add token and user info
}

type UserFavoritesResponse struct {
	// TODO: Add favorites list
}

type UserReviewsResponse struct {
	// TODO: Add reviews list
}