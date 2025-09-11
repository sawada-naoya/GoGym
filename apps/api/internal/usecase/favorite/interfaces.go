package favorite

import (
	"context"
)

// FavoriteUseCase defines the interface for favorite business logic
type FavoriteUseCase interface {
	// TODO: Add gym to favorites
	CreateFavorite(ctx context.Context, userID int64, gymID int64) error
	// TODO: Remove gym from favorites
	DeleteFavorite(ctx context.Context, userID int64, gymID int64) error
	// TODO: Check if gym is favorited by user
	IsFavorite(ctx context.Context, userID int64, gymID int64) (bool, error)
	// TODO: Get user's favorites with gym details
	GetUserFavorites(ctx context.Context, userID int64) (*FavoritesResponse, error)
}

// Repository interface
type FavoriteRepository interface {
	// TODO: Repository methods
}

// Response DTOs
type FavoritesResponse struct {
	// TODO: Add favorites with gym info
}