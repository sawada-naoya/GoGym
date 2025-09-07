package gym

import (
	"context"
	"gogym-api/internal/domain/gym"
)

// AddFavorite adds a gym to user's favorites
func (uc *UseCase) AddFavorite(ctx context.Context, req FavoriteGymRequest) error {
	uc.logger.InfoContext(ctx, "adding favorite gym",
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
	_, err := uc.gymRepo.FindByID(ctx, req.GymID)
	if err != nil {
		uc.logger.ErrorContext(ctx, "gym not found", "gym_id", req.GymID, "error", err)
		return gym.NewDomainErrorWithCause(err, "gym_not_found", "gym not found")
	}

	// Add to favorites
	if err := uc.favoriteRepo.AddFavorite(ctx, req.UserID, req.GymID); err != nil {
		uc.logger.ErrorContext(ctx, "failed to add favorite", "error", err)
		return gym.NewDomainErrorWithCause(err, "favorite_add_failed", "failed to add favorite")
	}

	uc.logger.InfoContext(ctx, "favorite gym added successfully",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	return nil
}

// RemoveFavorite removes a gym from user's favorites
func (uc *UseCase) RemoveFavorite(ctx context.Context, req FavoriteGymRequest) error {
	uc.logger.InfoContext(ctx, "removing favorite gym",
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
	if err := uc.favoriteRepo.RemoveFavorite(ctx, req.UserID, req.GymID); err != nil {
		uc.logger.ErrorContext(ctx, "failed to remove favorite", "error", err)
		return gym.NewDomainErrorWithCause(err, "favorite_remove_failed", "failed to remove favorite")
	}

	uc.logger.InfoContext(ctx, "favorite gym removed successfully",
		"user_id", req.UserID,
		"gym_id", req.GymID,
	)

	return nil
}

// GetFavoriteGyms retrieves user's favorite gyms
func (uc *UseCase) GetFavoriteGyms(ctx context.Context, req GetFavoriteGymsRequest) (*GetFavoriteGymsResponse, error) {
	uc.logger.InfoContext(ctx, "getting favorite gyms",
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

	result, err := uc.favoriteRepo.GetFavoriteGyms(ctx, req.UserID, pagination)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to get favorite gyms", "error", err)
		return nil, gym.NewDomainErrorWithCause(err, "favorites_fetch_failed", "failed to get favorite gyms")
	}

	return &GetFavoriteGymsResponse{
		Gyms:       result.Items,
		NextCursor: &result.NextCursor,
		HasMore:    result.HasMore,
	}, nil
}

// IsFavorite checks if a gym is in user's favorites
func (uc *UseCase) IsFavorite(ctx context.Context, userID gym.ID, gymID gym.ID) (bool, error) {
	if userID == 0 || gymID == 0 {
		return false, gym.NewDomainError(gym.ErrInvalidInput, "invalid_ids", "user ID and gym ID are required")
	}

	return uc.favoriteRepo.IsFavorite(ctx, userID, gymID)
}