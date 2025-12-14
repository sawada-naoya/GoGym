package user

import (
	"context"

	dom "gogym-api/internal/domain/entities"

	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *dom.User) error {
	recordUser := FromEntity(user)

	result := r.db.WithContext(ctx).Create(recordUser)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id ulid.ULID) (*dom.User, error) {
	var recordUser User

	err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id.String()).
		First(&recordUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // user not found
		}
		return nil, err
	}

	return ToEntity(&recordUser)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*dom.User, error) {
	var recordUser User

	err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("email = ?", email).
		First(&recordUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // user not found
		}
		return nil, err
	}

	return ToEntity(&recordUser)
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error

	if result != nil {
		return false, result
	}

	return count > 0, nil
}
