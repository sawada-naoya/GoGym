package workout

import (
	"context"
	"errors"
	"fmt"
	"time"

	wu "gogym-api/internal/application/workout"
	domain "gogym-api/internal/domain/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) wu.Repository {
	return &workoutRepository{db: db}
}

// insertWorkoutSets は WorkoutSet を1件ずつインサートするヘルパー関数
// GORM の Create() はリレーション（Sets）も自動作成するため、
// ID の重複を防ぐために Sets を nil にした後、手動で1件ずつ作成する
func (r *workoutRepository) insertWorkoutSets(tx *gorm.DB, recordID int, sets []WorkoutSet) error {
	for i := range sets {
		// 種目IDが無効な場合はスキップ（未選択の種目を保存しない）
		if sets[i].WorkoutExerciseID <= 0 {
			continue
		}

		newSet := WorkoutSet{
			WorkoutRecordID:   recordID,
			WorkoutExerciseID: sets[i].WorkoutExerciseID,
			SetNumber:         sets[i].SetNumber,
			WeightKg:          sets[i].WeightKg,
			Reps:              sets[i].Reps,
			EstimatedMax:      sets[i].EstimatedMax,
			Note:              sets[i].Note,
		}
		if err := tx.Create(&newSet).Error; err != nil {
			return fmt.Errorf("failed to insert workout set: %w", err)
		}
	}
	return nil
}

// GetRecordsByDate は指定日付のワークアウトレコードを取得（全部位）
// レコードが存在しない場合は空のドメインモデルを返す
func (r *workoutRepository) GetRecordsByDate(ctx context.Context, userID string, date time.Time) (domain.WorkoutRecord, error) {
	var record WorkoutRecord
	err := r.db.WithContext(ctx).
		Preload("Gym").
		Preload("Sets", func(db *gorm.DB) *gorm.DB {
			return db.Order("workout_sets.set_number ASC")
		}).
		Preload("Sets.Exercise").
		Preload("Sets.Exercise.Part").
		Preload("Sets.Exercise.Part.Translations").
		Where("user_id = ? AND performed_date = ?", userID, date).
		First(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// レコードが存在しない場合は空のドメインエンティティを返す
			return domain.WorkoutRecord{}, nil
		}
		return domain.WorkoutRecord{}, fmt.Errorf("error fetching workout records: %w", err)
	}

	// リポジトリモデルをドメインエンティティに変換
	domainRecord := ToEntity(&record)
	if domainRecord == nil {
		return domain.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain entity")
	}

	return *domainRecord, nil
}

// CreateWorkoutRecord は新規ワークアウトレコードを作成
// トランザクション内で Record と Sets を別々に作成し、ID の重複を防ぐ
func (r *workoutRepository) CreateWorkoutRecord(ctx context.Context, workout domain.WorkoutRecord) error {
	recordWorkout := FromEntity(&workout)
	if recordWorkout == nil {
		return fmt.Errorf("failed to convert domain workout record to repository record")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Sets を一時保存し、Record のみ作成（Sets の自動作成を防ぐ）
		sets := recordWorkout.Sets
		recordWorkout.Sets = nil

		if err := tx.Create(recordWorkout).Error; err != nil {
			return fmt.Errorf("failed to create workout record: %w", err)
		}

		// Sets を手動で作成
		return r.insertWorkoutSets(tx, recordWorkout.ID, sets)
	})
}

// UpsertWorkoutRecord は同日のレコードがあれば更新、なければ新規作成
// - 同じ日付のレコードが存在: メタデータを更新し、セットを追加/置き換え
// - 存在しない: 新規作成
func (r *workoutRepository) UpsertWorkoutRecord(ctx context.Context, workout domain.WorkoutRecord) error {
	recordWorkout := FromEntity(&workout)
	if recordWorkout == nil {
		return fmt.Errorf("failed to convert domain workout record to repository record")
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 同日のレコードを検索（部位に関係なく）
		var existingRecord WorkoutRecord
		err := tx.
			Preload("Sets.Exercise").
			Where("user_id = ? AND performed_date = ?", recordWorkout.UserID, recordWorkout.PerformedDate).
			First(&existingRecord).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新規作成パス
			return r.createRecordWithSets(tx, recordWorkout)
		}
		if err != nil {
			return fmt.Errorf("failed to check existing record: %w", err)
		}

		// 更新パス：メタデータを更新し、新しいセットを追加
		// 最初のセットから部位IDを取得（既存の同部位セットを削除するため）
		var partID *int
		if len(recordWorkout.Sets) > 0 {
			var exercise WorkoutExercise
			if err := tx.First(&exercise, recordWorkout.Sets[0].WorkoutExerciseID).Error; err == nil {
				partID = exercise.WorkoutPartID
			}
		}

		return r.updateRecordAndReplaceSets(tx, &existingRecord, recordWorkout, partID)
	})
}

// createRecordWithSets は新規レコードとセットを作成
func (r *workoutRepository) createRecordWithSets(tx *gorm.DB, recordWorkout *WorkoutRecord) error {
	sets := recordWorkout.Sets
	recordWorkout.Sets = nil
	recordWorkout.ID = 0 // 新規作成なのでIDをクリア（オートインクリメント）

	if err := tx.Create(recordWorkout).Error; err != nil {
		return fmt.Errorf("failed to create workout record: %w", err)
	}

	return r.insertWorkoutSets(tx, recordWorkout.ID, sets)
}

// updateRecordAndReplaceSets は既存レコードのメタデータを更新し、セットを置き換え
func (r *workoutRepository) updateRecordAndReplaceSets(tx *gorm.DB, existing *WorkoutRecord, new *WorkoutRecord, partID *int) error {
	// メタデータ（時刻・コンディション・ノート・ジムID）を更新
	updates := map[string]interface{}{
		"started_at":      new.StartedAt,
		"ended_at":        new.EndedAt,
		"note":            new.Note,
		"condition_level": new.ConditionLevel,
		"gym_id":          new.GymID,
	}
	if err := tx.Model(existing).Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update workout record metadata: %w", err)
	}

	// 同じ部位のセットを削除
	if err := r.deleteSetsByPart(tx, existing.ID, partID); err != nil {
		return fmt.Errorf("failed to delete existing sets: %w", err)
	}

	// 無効なセット（exerciseID=0や存在しない種目）を削除
	if err := r.deleteInvalidSets(tx, existing.ID); err != nil {
		return fmt.Errorf("failed to delete invalid sets: %w", err)
	}

	// 新しいセットを作成
	return r.insertWorkoutSets(tx, existing.ID, new.Sets)
}

// deleteSetsByPart は指定部位のセットを物理削除
// 物理削除を使用してユニーク制約の問題を回避
func (r *workoutRepository) deleteSetsByPart(tx *gorm.DB, recordID int, partID *int) error {
	if partID == nil {
		return nil
	}

	// 指定部位に属する種目IDを取得
	var exerciseIDs []int
	if err := tx.Model(&WorkoutExercise{}).
		Select("id").
		Where("workout_part_id = ?", *partID).
		Pluck("id", &exerciseIDs).Error; err != nil {
		return fmt.Errorf("failed to get exercise IDs: %w", err)
	}

	if len(exerciseIDs) == 0 {
		return nil
	}

	// 該当するセットを物理削除
	if err := tx.Unscoped().
		Where("workout_record_id = ?", recordID).
		Where("workout_exercise_id IN ?", exerciseIDs).
		Delete(&WorkoutSet{}).Error; err != nil {
		return fmt.Errorf("failed to delete sets by part: %w", err)
	}

	return nil
}

// deleteInvalidSets は無効なセット（weight=0 かつ reps=0、またはexerciseID=0）を物理削除
func (r *workoutRepository) deleteInvalidSets(tx *gorm.DB, recordID int) error {
	// weight_kg=0 かつ reps=0 のセット、またはexerciseID=0のセットを削除
	if err := tx.Unscoped().
		Where("workout_record_id = ?", recordID).
		Where("(weight_kg = 0 AND reps = 0) OR workout_exercise_id = 0 OR workout_exercise_id IS NULL").
		Delete(&WorkoutSet{}).Error; err != nil {
		return fmt.Errorf("failed to delete invalid sets: %w", err)
	}

	return nil
}

// GetWorkoutParts はユーザーのワークアウト部位一覧を取得
// 各部位に紐づく種目と翻訳データもプリロード
func (r *workoutRepository) GetWorkoutParts(ctx context.Context, userID string) ([]dom.WorkoutPart, error) {
	var parts []WorkoutPart

	err := r.db.WithContext(ctx).
		Preload("Translations").
		Preload("Exercises", "user_id = ?", userID).
		Where("user_id = ?", userID).
		Order("key ASC").
		Find(&parts).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching workout parts: %w", err)
	}

	return WorkoutPartsToDomain(parts), nil
}

// CountUserWorkoutParts はユーザーのワークアウト部位数をカウント
func (r *workoutRepository) CountUserWorkoutParts(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&WorkoutPart{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("error counting workout parts: %w", err)
	}

	return count, nil
}

// CreateWorkoutParts は複数のワークアウト部位を一括作成（翻訳データも含む）
func (r *workoutRepository) CreateWorkoutParts(ctx context.Context, userID string, parts []domain.WorkoutPart) error {
	recordParts := make([]WorkoutPart, len(parts))
	for i, part := range parts {
		translations := make([]WorkoutPartTranslation, len(part.Translations))
		for j, trans := range part.Translations {
			translations[j] = WorkoutPartTranslation{
				Locale: trans.Locale,
				Name:   trans.Name,
			}
		}

		recordParts[i] = WorkoutPart{
			Key:          part.Key,
			UserID:       &userID,
			Translations: translations,
		}
	}

	err := r.db.WithContext(ctx).Create(&recordParts).Error
	if err != nil {
		return fmt.Errorf("error creating workout parts: %w", err)
	}

	return nil
}

// DeleteWorkoutExercise はワークアウト種目を論理削除
func (r *workoutRepository) DeleteWorkoutExercise(ctx context.Context, userID string, exerciseID int64) error {
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", exerciseID, userID).
		Delete(&WorkoutExercise{}).Error
	if err != nil {
		return fmt.Errorf("error deleting workout exercise: %w", err)
	}

	return nil
}

// UpsertWorkoutExercises はワークアウト種目を一括 upsert
// - ID が指定されていれば更新、なければ新規作成
// - OnConflict で ID 衝突時は name と workout_part_id を更新
func (r *workoutRepository) UpsertWorkoutExercises(ctx context.Context, userID string, exercises []domain.WorkoutExerciseRef) error {
	recordExercises := make([]WorkoutExercise, 0, len(exercises))
	for _, exercise := range exercises {
		var partID *int
		if exercise.PartID != nil {
			pid := int(*exercise.PartID)
			partID = &pid
		}

		recordExercise := WorkoutExercise{
			Name:          exercise.Name,
			WorkoutPartID: partID,
			UserID:        &userID,
		}

		if exercise.ID != 0 {
			recordExercise.ID = int(exercise.ID)
		}

		recordExercises = append(recordExercises, recordExercise)
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"name", "workout_part_id"}),
		}).
		Create(&recordExercises).Error
	if err != nil {
		return fmt.Errorf("error upserting workout exercises: %w", err)
	}

	return nil
}

func (r *workoutRepository) GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (domain.WorkoutRecord, error) {
	var rec WorkoutRecord

	// サブクエリ: 指定したエクササイズを含む最新のレコードIDを取得
	// performed_date（実施日）とid（より新しいレコード）で最新を判定
	subQuery := r.db.Table("workout_records").
		Select("workout_records.id").
		Joins("INNER JOIN workout_sets ON workout_sets.workout_record_id = workout_records.id").
		Where("workout_records.user_id = ? AND workout_sets.workout_exercise_id = ?", userID, exerciseID).
		Order("workout_records.performed_date DESC, workout_records.id DESC").
		Limit(1)

	err := r.db.WithContext(ctx).
		Preload("Gym").
		Preload("Sets", "workout_exercise_id = ?", exerciseID). // 指定エクササイズのセットのみ取得
		Preload("Sets.Exercise").
		Preload("Sets.Exercise.Part").
		Where("user_id = ? AND id = (?)", userID, subQuery).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.WorkoutRecord{}, nil
		}
		return domain.WorkoutRecord{}, fmt.Errorf("error fetching last workout record: %w", err)
	}

	domainRecord := ToEntity(&rec)
	if domainRecord == nil {
		return domain.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain")
	}

	return *domainRecord, nil
}
