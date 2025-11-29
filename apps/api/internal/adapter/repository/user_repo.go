package repository

import (
	"context"
	"gogym-api/internal/adapter/repository/mapper"
	"gogym-api/internal/adapter/repository/record"
	dom "gogym-api/internal/domain/user"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserRepository struct { // ← 公開にする
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *dom.User) error {
	recordUser := mapper.FromUserEntity(user)

	result := r.db.WithContext(ctx).Create(recordUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id ulid.ULID) (*dom.User, error) {
	var recordUser record.User

	err := r.db.WithContext(ctx).
		Model(&record.User{}).
		Where("id = ?", id.String()).
		First(&recordUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dom.NewDomainError("user_not_found")
		}
		return nil, err
	}

	return mapper.ToUserEntity(&recordUser), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email dom.Email) (*dom.User, error) {
	var recordUser record.User

	err := r.db.WithContext(ctx).
		Model(&record.User{}).
		Where("email = ?", email.String()).
		First(&recordUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dom.NewDomainError("user_not_found")
		}
		return nil, err
	}

	return mapper.ToUserEntity(&recordUser), nil
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email dom.Email) (bool, error) {
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
