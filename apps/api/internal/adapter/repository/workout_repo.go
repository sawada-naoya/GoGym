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

func (r *workoutRepository) Create(ctx context.Context, workout dom.WorkoutRecord) error {
	recordWorkout := mapper.WorkoutRecordToRecord(&workout)
	if recordWorkout == nil {
		return fmt.Errorf("failed to convert domain workout record to repository record")
	}

	// トランザクション内でWorkoutRecordとWorkoutSetsを作成
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// WorkoutRecordを作成
		if err := tx.Create(recordWorkout).Error; err != nil {
			return fmt.Errorf("failed to create workout record: %w", err)
		}

		// WorkoutSetsのWorkoutRecordIDを更新
		for i := range recordWorkout.Sets {
			recordWorkout.Sets[i].WorkoutRecordID = recordWorkout.ID
		}

		// WorkoutSetsを一括作成
		if len(recordWorkout.Sets) > 0 {
			if err := tx.Create(&recordWorkout.Sets).Error; err != nil {
				return fmt.Errorf("failed to create workout sets: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *workoutRepository) GetWorkoutParts(ctx context.Context, userID string) ([]dom.WorkoutPart, error) {
	var parts []record.WorkoutPart

	// ユーザー固有の部位のみを取得
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&parts).Error

	if err != nil {
		return nil, fmt.Errorf("error fetching workout parts: %w", err)
	}

	return mapper.WorkoutPartsToDomain(parts), nil
}

func (r *workoutRepository) CountUserWorkoutParts(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&record.WorkoutPart{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("error counting workout parts: %w", err)
	}

	return count, nil
}

func (r *workoutRepository) CreateWorkoutParts(ctx context.Context, userID string, parts []dom.WorkoutPart) error {
	recordParts := make([]record.WorkoutPart, len(parts))
	for i, part := range parts {
		recordParts[i] = record.WorkoutPart{
			Name:   part.Name,
			UserID: &userID,
		}
	}

	err := r.db.WithContext(ctx).Create(&recordParts).Error
	if err != nil {
		return fmt.Errorf("error creating workout parts: %w", err)
	}

	return nil
}

func (r *workoutRepository) CreateWorkoutExercises(ctx context.Context, userID string, exercises []dom.WorkoutExerciseRef) error {
	recordExercises := make([]record.WorkoutExercise, len(exercises))
	for i, exercise := range exercises {
		var partID *int
		if exercise.PartID != nil {
			pid := int(*exercise.PartID)
			partID = &pid
		}

		recordExercises[i] = record.WorkoutExercise{
			Name:          exercise.Name,
			WorkoutPartID: partID,
			UserID:        &userID,
		}
	}

	err := r.db.WithContext(ctx).Create(&recordExercises).Error
	if err != nil {
		return fmt.Errorf("error creating workout exercises: %w", err)
	}

	return nil
}
