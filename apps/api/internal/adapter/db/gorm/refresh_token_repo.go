package gorm

import (
	"context"

	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) Create(ctx context.Context) error {
	// TODO: Create refresh token
	return nil
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) error {
	// TODO: Find refresh token
	return nil
}

func (r *RefreshTokenRepository) Delete(ctx context.Context, token string) error {
	// TODO: Delete refresh token
	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	// TODO: Delete all user's refresh tokens
	return nil
}