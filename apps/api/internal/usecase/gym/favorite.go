package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// AddFavorite adds a gym to user's favorites
func (gu *GymUseCase) AddFavorite(ctx context.Context, req FavoriteGymRequest) error {
	gu.logger.InfoContext(ctx, "adding favorite gym",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	if req.UserID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if req.GymID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	// Check if gym exists
	_, err := gu.gymRepo.FindByID(ctx, req.GymID)
	if err != nil {
		gu.logger.ErrorContext(ctx, "gym not found", "gym_id", req.GymID, "error", err)
		return gym.NewDomainErrorWithCause(err, "gym_not_found", "gym not found")
	}

	// Add to favorites
	if err := gu.favoriteRepo.AddFavorite(ctx, req.UserID, req.GymID); err != nil {
		gu.logger.ErrorContext(ctx, "failed to add favorite", "error", err)
		return gym.NewDomainErrorWithCause(err, "favorite_add_failed", "failed to add favorite")
	}

	gu.logger.InfoContext(ctx, "favorite gym added successfully",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	return nil
}

// RemoveFavorite removes a gym from user's favorites
func (gu *GymUseCase) RemoveFavorite(ctx context.Context, req FavoriteGymRequest) error {
	gu.logger.InfoContext(ctx, "removing favorite gym",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	if req.UserID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	if req.GymID == 0 {
		return gym.NewDomainError(gym.ErrInvalidInput, "invalid_gym_id", "gym ID is required")
	}

	// Remove from favorites
	if err := gu.favoriteRepo.RemoveFavorite(ctx, req.UserID, req.GymID); err != nil {
		gu.logger.ErrorContext(ctx, "failed to remove favorite", "error", err)
		return gym.NewDomainErrorWithCause(err, "favorite_remove_failed", "failed to remove favorite")
	}

	gu.logger.InfoContext(ctx, "favorite gym removed successfully",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	return nil
}

// GetFavoriteGyms retrieves user's favorite gyms
func (gu *GymUseCase) GetFavoriteGyms(ctx context.Context, req GetFavoriteGymsRequest) (*GetFavoriteGymsResponse, error) {
	gu.logger.InfoContext(ctx, "getting favorite gyms",
		"user_id", req.UserID,
		"limit", req.Limit,
	)

	if req.UserID == 0 {
		return nil, gym.NewDomainError(gym.ErrInvalidInput, "invalid_user_id", "user ID is required")
	}

	// Set default limit
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}

	pagination := gym.Pagination{
		Cursor: req.Cursor,
		Limit:  req.Limit,
	}

	result, err := gu.favoriteRepo.GetFavoriteGyms(ctx, req.UserID, pagination)
	if err != nil {
		gu.logger.ErrorContext(ctx, "failed to get favorite gyms", "error", err)
		return nil, gym.NewDomainErrorWithCause(err, "favorites_fetch_failed", "failed to get favorite gyms")
	}

	return &GetFavoriteGymsResponse{
		Gyms:       result.Items,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// IsFavorite checks if a gym is in user's favorites
func (gu *GymUseCase) IsFavorite(ctx context.Context, userID gym.ID, gymID gym.ID) (bool, error) {
	if userID == 0 || gymID == 0 {
		return false, gym.NewDomainError(gym.ErrInvalidInput, "invalid_ids", "user ID and gym ID are required")
	}

	return gu.favoriteRepo.IsFavorite(ctx, userID, gymID)
}