package gym

import (
	"context"
	"errors"

	gu "gogym-api/internal/application/gym"
	domain "gogym-api/internal/domain/entities"

	"gorm.io/gorm"
)

type gymRepository struct {
	db *gorm.DB
}

func NewGymRepository(db *gorm.DB) gu.Repository {
	return &gymRepository{db: db}
}

// FindByNormalizedName finds a gym by normalized name and creator
func (r *gymRepository) FindByNormalizedName(ctx context.Context, createdBy string, normalizedName string) (*domain.Gym, error) {
	var record GymRecord
	err := r.db.WithContext(ctx).
		Where("created_by = ? AND normalized_name = ? AND deleted_at IS NULL", createdBy, normalizedName).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gu.ErrNotFound
		}
		return nil, err
	}

	return ToEntity(&record), nil
}

// CreateGym creates a new gym with duplicate key handling
func (r *gymRepository) CreateGym(ctx context.Context, createdBy string, name string, normalizedName string) (*domain.Gym, error) {
	record := &GymRecord{
		Name:           name,
		NormalizedName: normalizedName,
		CreatedBy:      createdBy,
		// Required fields with dummy values (actual gyms would have real coordinates)
		Latitude:  0,
		Longitude: 0,
		SourceURL: "",
	}

	err := r.db.WithContext(ctx).Create(record).Error
	if err != nil {
		// On duplicate key error, try to find and return the existing record
		if isDuplicateKeyError(err) {
			return r.FindByNormalizedName(ctx, createdBy, normalizedName)
		}
		return nil, err
	}

	return ToEntity(record), nil
}

// isDuplicateKeyError checks if the error is a duplicate key constraint violation
func isDuplicateKeyError(err error) bool {
	// PostgreSQL duplicate key error code: 23505
	// MySQL duplicate key error code: 1062
	// This is a simplified check - in production you might want to use a more robust method
	if err == nil {
		return false
	}
	errStr := err.Error()
	// PostgreSQL
	if contains(errStr, "duplicate key") || contains(errStr, "23505") {
		return true
	}
	// MySQL
	if contains(errStr, "Duplicate entry") || contains(errStr, "1062") {
		return true
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
