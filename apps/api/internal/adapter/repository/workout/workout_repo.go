package workout

import (
	"context"
	"errors"
	"fmt"
	"time"

	dom "gogym-api/internal/domain/entities"
	workoutUsecase "gogym-api/internal/usecase/workout"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type workoutRepository struct {
	db *gorm.DB
}

func NewWorkoutRepository(db *gorm.DB) workoutUsecase.Repository {
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
func (r *workoutRepository) GetRecordsByDate(ctx context.Context, userID string, date string) (dom.WorkoutRecord, error) {
	var rec WorkoutRecord
	err := r.db.WithContext(ctx).
		Preload("Gym").
		Preload("Sets").
		Preload("Sets.Exercise").
		Preload("Sets.Exercise.Part").
		Where("user_id = ? AND performed_date = ?", userID, date).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dom.WorkoutRecord{}, nil
		}
		return dom.WorkoutRecord{}, fmt.Errorf("error fetching workout records: %w", err)
	}

	domainRecord := ToEntity(&rec)
	if domainRecord == nil {
		return dom.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain")
	}

	return *domainRecord, nil
}

// GetRecordsByDateAndPart は指定日付と部位IDでワークアウトレコードを取得
// partID が nil の場合は全ての部位を取得、指定された場合は該当部位のみフィルタ
func (r *workoutRepository) GetRecordsByDateAndPart(ctx context.Context, userID string, date string, partID *int64) (dom.WorkoutRecord, error) {
	if partID == nil {
		return r.GetRecordsByDate(ctx, userID, date)
	}

	// 指定部位の種目のみを含むセットを Preload でフィルタリング
	var rec WorkoutRecord
	err := r.db.WithContext(ctx).
		Preload("Gym").
		Preload("Sets.Exercise.Part", "workout_parts.id = ?", *partID).
		Preload("Sets", "workout_exercise_id IN (SELECT id FROM workout_exercises WHERE workout_part_id = ?)", *partID).
		Preload("Sets.Exercise", "workout_part_id = ?", *partID).
		Where("user_id = ? AND performed_date = ?", userID, date).
		First(&rec).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dom.WorkoutRecord{}, nil
		}
		return dom.WorkoutRecord{}, fmt.Errorf("error fetching workout records: %w", err)
	}

	domainRecord := ToEntity(&rec)
	if domainRecord == nil {
		return dom.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain")
	}

	return *domainRecord, nil
}

// CreateWorkoutRecord は新規ワークアウトレコードを作成
// トランザクション内で Record と Sets を別々に作成し、ID の重複を防ぐ
func (r *workoutRepository) CreateWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
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

// UpsertWorkoutRecord は同日同部位のレコードがあれば更新、なければ新規作成
// - 同じ日付 & 同じ部位: メタデータ（時刻・場所・コンディション）を更新し、セットを置き換え
// - 異なる部位: 新規作成
func (r *workoutRepository) UpsertWorkoutRecord(ctx context.Context, workout dom.WorkoutRecord) error {
	recordWorkout := FromEntity(&workout)
	if recordWorkout == nil {
		return fmt.Errorf("failed to convert domain workout record to repository record")
	}

	// 最初のセットから部位IDを取得（同日同部位の検索に使用）
	partID, err := r.getPartIDFromFirstSet(ctx, recordWorkout.Sets)
	if err != nil {
		return fmt.Errorf("failed to get part ID: %w", err)
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existingRecord, err := r.findExistingRecord(tx, recordWorkout.UserID, recordWorkout.PerformedDate, partID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check existing record: %w", err)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 新規作成パス
			return r.createRecordWithSets(tx, recordWorkout)
		}

		// 更新パス
		return r.updateRecordAndReplaceSets(tx, existingRecord, recordWorkout, partID)
	})
}

// getPartIDFromFirstSet は最初のセットから部位IDを取得
func (r *workoutRepository) getPartIDFromFirstSet(ctx context.Context, sets []WorkoutSet) (*int, error) {
	if len(sets) == 0 {
		// セットが空の場合は部位を特定できない（全レコード検索になってしまう）
		// エラーを返すか、新規作成を強制する
		return nil, nil
	}

	var exercise WorkoutExercise
	if err := r.db.WithContext(ctx).First(&exercise, sets[0].WorkoutExerciseID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exercise: %w", err)
	}

	// 部位IDがnullの場合も問題（曖昧な検索になる）
	if exercise.WorkoutPartID == nil {
		return nil, fmt.Errorf("exercise has no workout_part_id")
	}

	return exercise.WorkoutPartID, nil
}

// findExistingRecord は同日同部位のレコードを検索
func (r *workoutRepository) findExistingRecord(tx *gorm.DB, userID string, performedDate time.Time, partID *int) (*WorkoutRecord, error) {
	var existingRecord WorkoutRecord

	query := tx.
		Preload("Sets.Exercise").
		Where("workout_records.user_id = ? AND workout_records.performed_date = ?", userID, performedDate)

	// 部位IDが指定されている場合のみフィルタ
	if partID != nil {
		query = query.
			Joins("JOIN workout_sets ON workout_sets.workout_record_id = workout_records.id").
			Joins("JOIN workout_exercises ON workout_exercises.id = workout_sets.workout_exercise_id").
			Where("workout_exercises.workout_part_id = ?", *partID)
	}

	err := query.First(&existingRecord).Error
	if err != nil {
		return nil, err
	}

	return &existingRecord, nil
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
func (r *workoutRepository) CreateWorkoutParts(ctx context.Context, userID string, parts []dom.WorkoutPart) error {
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
func (r *workoutRepository) UpsertWorkoutExercises(ctx context.Context, userID string, exercises []dom.WorkoutExerciseRef) error {
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

func (r *workoutRepository) GetLastWorkoutRecord(ctx context.Context, userID string, exerciseID int64) (dom.WorkoutRecord, error) {
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
			return dom.WorkoutRecord{}, nil
		}
		return dom.WorkoutRecord{}, fmt.Errorf("error fetching last workout record: %w", err)
	}

	domainRecord := ToEntity(&rec)
	if domainRecord == nil {
		return dom.WorkoutRecord{}, fmt.Errorf("failed to convert record to domain")
	}

	return *domainRecord, nil
}
