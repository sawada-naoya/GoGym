// internal/adapter/db/gorm/user_repo.go
// 役割: ユーザードメイン用リポジトリ実装（Infrastructure Layer）
// record↔entity変換を行いDB操作する実装。ドメインエンティティとGORMレコード間の変換を担当
package gorm

import (
	"context"
	"gogym-api/internal/adapter/db/gorm/record"
	dom "gogym-api/internal/domain/user"
	uc "gogym-api/internal/usecase/user"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository は新しいユーザーリポジトリを作成する
func NewUserRepository(db *gorm.DB) uc.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *dom.User) error {
	// ドメインエンティティをGORMレコードに変換
	recordUser := FromUserEntity(user)

	result := r.db.WithContext(ctx).Create(recordUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email dom.Email) (*dom.User, error) {
	var recordUser record.User

	err := r.db.WithContext(ctx).
		Model(&record.User{}).
		Where("email = ?", email).
		First(&recordUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dom.NewDomainError("user_not_found")
		}
		return nil, err
	}

	return ToUserEntity(&recordUser), nil
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email dom.Email) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&record.User{}).
		Where("email = ?", email.String()).
		Count(&count).Error

	if result != nil {
		return false, result
	}

	return count > 0, nil
}
