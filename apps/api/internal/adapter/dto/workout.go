package dto

import (
	"fmt"
	"time"

	dom "gogym-api/internal/domain/workout"
	"gogym-api/internal/util"
)

// WorkoutPartDTO represents a workout part (e.g., chest, back, legs)
type WorkoutPartListItemDTO struct {
	ID        int64                         `json:"id"`
	Name      string                        `json:"name"`
	Exercises []WorkoutExerciseListItemDTO `json:"exercises"`
}

type WorkoutRecordDTO struct {
	ID             *int64  `json:"id,omitempty"`         // 既存ならrecord id
	PerformedDate  string  `json:"performed_date"`       // "YYYY-MM-DD"
	StartedAt      *string `json:"started_at,omitempty"` // "HH:mm"
	EndedAt        *string `json:"ended_at,omitempty"`   // "HH:mm"
	Place          string  `json:"place"`
	Note           *string `json:"note,omitempty"`
	ConditionLevel *int    `json:"condition_level,omitempty"` // 1..5

	WorkoutPart WorkoutPartDTO `json:"workout_part"`
	Exercises   []ExerciseDTO  `json:"exercises"`
}

type WorkoutPartDTO struct {
	ID     *int64  `json:"id,omitempty"`
	Name   *string `json:"name,omitempty"`
	Source *string `json:"source,omitempty"` // "preset" | "custom" | null
}

type ExerciseDTO struct {
	ID            *int64   `json:"id,omitempty"`
	Name          string   `json:"name"` // 種目名
	WorkoutPartID *int64   `json:"workout_part_id,omitempty"`
	Sets          []SetDTO `json:"sets"`
}

type SetDTO struct {
	ID        *int64   `json:"id,omitempty"`
	SetNumber int      `json:"set_number"`
	WeightKg  *float64 `json:"weight_kg,omitempty"` // 空文字→null→nil→層内で検証
	Reps      *int     `json:"reps,omitempty"`
	Note      *string  `json:"note,omitempty"`
}

type CreateWorkoutExerciseRequest struct {
	Exercises []CreateWorkoutExerciseItem `json:"exercises"`
}

type CreateWorkoutExerciseItem struct {
	ID            *int64 `json:"id,omitempty"` // nil = insert, value = update
	Name          string `json:"name"`
	WorkoutPartID int64  `json:"workout_part_id"`
}

type WorkoutExerciseListItemDTO struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	WorkoutPartID *int64 `json:"workout_part_id,omitempty"`
}

// DomainToDTO converts domain.WorkoutRecord to WorkoutFormDTO
func WorkoutRecordToDTO(record *dom.WorkoutRecord) *WorkoutRecordDTO {
	if record == nil {
		return nil
	}

	dto := &WorkoutRecordDTO{
		ID:             domainIDToInt64Ptr(record.ID),
		PerformedDate:  util.FormatJSTDate(record.PerformedDate),
		StartedAt:      timeToJSTHHmm(record.StartedAt),
		EndedAt:        timeToJSTHHmm(record.EndedAt),
		Place:          stringPtrToString(record.Place),
		Note:           record.Note,
		ConditionLevel: conditionLevelToIntPtr(record.Condition),
		WorkoutPart:    WorkoutPartDTO{}, // WorkoutPart情報がない場合は空
		Exercises:      []ExerciseDTO{},
	}

	// Setsを Exercise ごとにグループ化
	exerciseMap := make(map[dom.ID]*ExerciseDTO)
	for _, set := range record.Sets {
		exerciseID := set.Exercise.ID

		// 初めて見るExerciseの場合、ExerciseDTOを作成
		if _, exists := exerciseMap[exerciseID]; !exists {
			exerciseMap[exerciseID] = &ExerciseDTO{
				ID:            domainIDToInt64Ptr(&exerciseID),
				Name:          set.Exercise.Name,
				WorkoutPartID: domainIDToInt64Ptr(set.Exercise.PartID),
				Sets:          []SetDTO{},
			}
		}

		// SetDTOを追加
		weightKg := float64(set.Weight)
		reps := int(set.Reps)
		exerciseMap[exerciseID].Sets = append(exerciseMap[exerciseID].Sets, SetDTO{
			ID:        domainIDToInt64Ptr(set.ID),
			SetNumber: set.SetNumber,
			WeightKg:  &weightKg,
			Reps:      &reps,
			Note:      set.Note,
		})
	}

	// Mapから配列に変換
	for _, exercise := range exerciseMap {
		dto.Exercises = append(dto.Exercises, *exercise)
	}

	return dto
}

// Helper functions
func domainIDToInt64Ptr(id *dom.ID) *int64 {
	if id == nil {
		return nil
	}
	i := int64(*id)
	return &i
}

func timeToJSTHHmm(t *time.Time) *string {
	if t == nil {
		return nil
	}
	// UTC → JST 変換してから HH:mm 形式で返す
	jstTime := util.ToJST(*t)
	hhmm := fmt.Sprintf("%02d:%02d", jstTime.Hour(), jstTime.Minute())
	return &hhmm
}

func stringPtrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func conditionLevelToIntPtr(c dom.ConditionLevel) *int {
	if c == dom.CondUnknown {
		return nil
	}
	i := int(c)
	return &i
}

// WorkoutRecordDTOToDomain converts WorkoutRecordDTO to domain.WorkoutRecord
func WorkoutRecordDTOToDomain(dto *WorkoutRecordDTO) (*dom.WorkoutRecord, error) {
	if dto == nil {
		return nil, fmt.Errorf("dto is nil")
	}

	// Parse performedDate (JST として解釈)
	jstLoc, _ := time.LoadLocation("Asia/Tokyo")
	performedDate, err := time.ParseInLocation("2006-01-02", dto.PerformedDate, jstLoc)
	if err != nil {
		return nil, fmt.Errorf("invalid performedDate format: %w", err)
	}

	// UTC に変換してから WorkoutRecord を作成
	performedDateUTC := performedDate.UTC()

	// Create WorkoutRecord
	record, err := dom.NewWorkoutRecord(dom.ULID(""), performedDateUTC) // userID will be set by handler/usecase
	if err != nil {
		return nil, fmt.Errorf("failed to create workout record: %w", err)
	}

	// Set ID if exists
	if dto.ID != nil {
		id := dom.ID(*dto.ID)
		record.ID = &id
	}

	// Parse and set times
	var startedAt, endedAt *time.Time
	if dto.StartedAt != nil {
		t, err := parseTimeWithDate(performedDate, *dto.StartedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid startedAt format: %w", err)
		}
		startedAt = &t
	}
	if dto.EndedAt != nil {
		t, err := parseTimeWithDate(performedDate, *dto.EndedAt)
		if err != nil {
			return nil, fmt.Errorf("invalid endedAt format: %w", err)
		}
		endedAt = &t
	}
	if err := record.SetTimes(startedAt, endedAt); err != nil {
		return nil, fmt.Errorf("invalid times: %w", err)
	}

	// Set other fields
	if dto.Place != "" {
		record.Place = &dto.Place
	}
	record.Note = dto.Note
	if dto.ConditionLevel != nil {
		record.Condition = dom.ConditionLevel(*dto.ConditionLevel)
	}

	// Convert exercises to sets
	for _, exercise := range dto.Exercises {
		exerciseRef := dom.WorkoutExerciseRef{
			Name: exercise.Name,
		}
		if exercise.ID != nil {
			exerciseRef.ID = dom.ID(*exercise.ID)
		}
		if exercise.WorkoutPartID != nil {
			partID := dom.ID(*exercise.WorkoutPartID)
			exerciseRef.PartID = &partID
		}
		// Owner は exercise.ID の有無で判定される想定
		// 既存の exercise.ID がない場合は新規作成されるため、Owner は usecase 層で設定される

		// Convert sets
		for _, setDTO := range exercise.Sets {
			workoutSet := dom.WorkoutSet{
				Exercise:  exerciseRef,
				SetNumber: setDTO.SetNumber,
			}

			if setDTO.ID != nil {
				id := dom.ID(*setDTO.ID)
				workoutSet.ID = &id
			}

			if setDTO.WeightKg != nil {
				workoutSet.Weight = dom.WeightKg(*setDTO.WeightKg)
			}
			if setDTO.Reps != nil {
				workoutSet.Reps = dom.Reps(*setDTO.Reps)
			}
			workoutSet.Note = setDTO.Note

			// Add set to record
			if err := record.AddSet(workoutSet); err != nil {
				return nil, fmt.Errorf("failed to add set: %w", err)
			}
		}
	}

	return record, nil
}

// parseTimeWithDate combines a date and HH:mm time string
// フロントエンドからJST形式で受け取った日付時刻をUTCに変換してDBに保存
func parseTimeWithDate(date time.Time, hhmmStr string) (time.Time, error) {
	t, err := time.Parse("15:04", hhmmStr)
	if err != nil {
		return time.Time{}, err
	}

	// JST として解釈
	jstLoc, _ := time.LoadLocation("Asia/Tokyo")
	jstTime := time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, jstLoc)

	// UTC に変換
	return jstTime.UTC(), nil
}

// WorkoutPartToDTO converts domain.WorkoutPart to WorkoutPartListItemDTO
func WorkoutPartToDTO(part *dom.WorkoutPart) *WorkoutPartListItemDTO {
	if part == nil {
		return nil
	}

	// Exercisesを変換
	exercises := make([]WorkoutExerciseListItemDTO, 0, len(part.Exercises))
	for _, ex := range part.Exercises {
		var partIDPtr *int64
		if ex.PartID != nil {
			pid := int64(*ex.PartID)
			partIDPtr = &pid
		}

		exercises = append(exercises, WorkoutExerciseListItemDTO{
			ID:            int64(ex.ID),
			Name:          ex.Name,
			WorkoutPartID: partIDPtr,
		})
	}

	return &WorkoutPartListItemDTO{
		ID:        int64(part.ID),
		Name:      part.Name,
		Exercises: exercises,
	}
}

// WorkoutPartsToDTO converts slice of domain.WorkoutPart to slice of WorkoutPartListItemDTO
func WorkoutPartsToDTO(parts []dom.WorkoutPart) []WorkoutPartListItemDTO {
	result := make([]WorkoutPartListItemDTO, len(parts))
	for i, part := range parts {
		result[i] = *WorkoutPartToDTO(&part)
	}
	return result
}
