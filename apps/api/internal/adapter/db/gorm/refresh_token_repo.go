// internal/adapter/db/gorm/refresh_token_repo.go
// 役割: リフレッシュトークンドメイン用リポジトリ実装（Infrastructure Layer）
// record↔entity変換を行いDB操作する実装。ドメインエンティティとGORMレコード間の変換を担当
package gorm

import (
	"context"
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/common"
	"gogym-api/internal/domain/user"
	userUsecase "gogym-api/internal/usecase/user"
	"gorm.io/gorm"
	"time"
)

// refreshTokenRepository はuser.RefreshTokenRepositoryインターフェースを実装する
type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository は新しいリフレッシュトークンリポジトリを作成する
func NewRefreshTokenRepository(db *gorm.DB) userUsecase.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create は新しいリフレッシュトークンを作成する
func (r *refreshTokenRepository) Create(ctx context.Context, tokenEntity *user.RefreshToken) error {
	tokenRecord := record.FromRefreshTokenEntity(tokenEntity)
	
	if err := r.db.WithContext(ctx).Create(tokenRecord).Error; err != nil {
		return err
	}

	// 生成されたIDでエンティティを更新
	tokenEntity.ID = common.ID(tokenRecord.ID)
	tokenEntity.CreatedAt = tokenRecord.CreatedAt

	return nil
}

// FindByTokenHash はトークンハッシュでリフレッシュトークンを検索する
func (r *refreshTokenRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*user.RefreshToken, error) {
	var tokenRecord record.RefreshTokenRecord
	if err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&tokenRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, common.NewDomainError(common.ErrNotFound, "token_not_found", "refresh token not found")
		}
		return nil, err
	}

	return tokenRecord.ToEntity(), nil
}

// DeleteByTokenHash はトークンハッシュでリフレッシュトークンを削除する
func (r *refreshTokenRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	result := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).Delete(&record.RefreshTokenRecord{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return common.NewDomainError(common.ErrNotFound, "token_not_found", "refresh token not found")
	}

	return nil
}

// DeleteExpiredTokens は期限切れのすべてのリフレッシュトークンを削除する
func (r *refreshTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&record.RefreshTokenRecord{}).Error
}

// DeleteAllByUserID は特定のユーザーのすべてのリフレッシュトークンを削除する
func (r *refreshTokenRepository) DeleteAllByUserID(ctx context.Context, userID common.ID) error {
	return r.db.WithContext(ctx).Where("user_id = ?", int64(userID)).Delete(&record.RefreshTokenRecord{}).Error
}