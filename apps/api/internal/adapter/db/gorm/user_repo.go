// internal/adapter/db/gorm/user_repo.go
// 役割: ユーザードメイン用リポジトリ実装（Infrastructure Layer）
// record↔entity変換を行いDB操作する実装。ドメインエンティティとGORMレコード間の変換を担当
package gorm

import (
	"context"
	"gogym-api/internal/adapter/db/gorm/record"
	"gogym-api/internal/domain/user"
	userUsecase "gogym-api/internal/usecase/user"
	"gorm.io/gorm"
)

// userRepository はuser.Repositoryインターフェースを実装する
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しいユーザーリポジトリを作成する
func NewUserRepository(db *gorm.DB) userUsecase.Repository {
	return &userRepository{db: db}
}

// FindByID はIDでユーザーを検索する
func (r *userRepository) FindByID(ctx context.Context, id user.ID) (*user.User, error) {
	var userRecord record.UserRecord
	if err := r.db.WithContext(ctx).First(&userRecord, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.NewDomainError(user.ErrNotFound, "user_not_found", "user not found")
		}
		return nil, err
	}

	return ToUserEntity(&userRecord)
}

// FindByEmail はメールアドレスでユーザーを検索する
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var userRecord record.UserRecord
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userRecord).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, user.NewDomainError(user.ErrNotFound, "user_not_found", "user not found")
		}
		return nil, err
	}

	return ToUserEntity(&userRecord)
}

// Create は新しいユーザーを作成する
func (r *userRepository) Create(ctx context.Context, userEntity *user.User) error {
	userRecord := FromUserEntity(userEntity)
	
	if err := r.db.WithContext(ctx).Create(userRecord).Error; err != nil {
		return err
	}

	// 生成されたIDでエンティティを更新
	userEntity.ID = user.ID(userRecord.ID)
	userEntity.CreatedAt = userRecord.CreatedAt
	userEntity.UpdatedAt = userRecord.UpdatedAt

	return nil
}

// Update は既存のユーザーを更新する
func (r *userRepository) Update(ctx context.Context, userEntity *user.User) error {
	userRecord := FromUserEntity(userEntity)
	
	if err := r.db.WithContext(ctx).Save(userRecord).Error; err != nil {
		return err
	}

	// 更新タイムスタンプでエンティティを更新
	userEntity.UpdatedAt = userRecord.UpdatedAt

	return nil
}

// Delete はIDでユーザーを削除する
func (r *userRepository) Delete(ctx context.Context, id user.ID) error {
	result := r.db.WithContext(ctx).Delete(&record.UserRecord{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return user.NewDomainError(user.ErrNotFound, "user_not_found", "user not found")
	}

	return nil
}