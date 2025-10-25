package repository

import (
	"context"
	"errors"
	"fmt"

	"gogym-api/internal/adapter/repository/mapper"
	"gogym-api/internal/adapter/repository/record"
	dom "gogym-api/internal/domain/workout"
	workoutUsecase "gogym-api/internal/usecase/workout"

	"gorm.io/gorm"
)

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) workoutUsecase.Repository {
	return &workoutRepository{db: db}
}

func (r *workoutRepository) GetRecordsByDate(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error) {
	var rec record.WorkoutRecord
	err := r.db.WithContext(ctx).
		Preload("Sets.Exercise.Part").
		Where("user_id = ? AND performed_date = ?", userID, date).
		First(&rec).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// レコードが見つからない場合は空のドメインモデルを返す
			return dom.WorkoutRecord{}, nil
		}
		return dom.WorkoutRecord{}, fmt.Errorf("error fetching workout records: %w", err)
	}

	domainRecord := mapper.WorkoutRecordToDomain(&rec)
	if domainRecord == nil {
		return dom.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain")
	}

	return *domainRecord, nil
}
