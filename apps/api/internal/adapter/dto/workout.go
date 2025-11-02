package dto

import (
	"fmt"
	"time"

	dom "gogym-api/internal/domain/workout"
)

type WorkoutRecordDTO struct {
	ID             *int64  `json:"id,omitempty"`        // 既存ならrecord id
	PerformedDate  string  `json:"performedDate"`       // "YYYY-MM-DD"
	StartedAt      *string `json:"startedAt,omitempty"` // "HH:mm"
	EndedAt        *string `json:"endedAt,omitempty"`   // "HH:mm"
	Place          string  `json:"place"`
	Note           *string `json:"note,omitempty"`
	ConditionLevel *int    `json:"conditionLevel,omitempty"` // 1..5

	WorkoutPart WorkoutPartDTO `json:"workoutPart"`
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
	WorkoutPartID *int64   `json:"workoutPartId,omitempty"`
	IsDefault     *bool    `json:"isDefault,omitempty"` // 0|1 は bool で受けて層内で変換
	Sets          []SetDTO `json:"sets"`
}

type SetDTO struct {
	ID        *int64   `json:"id,omitempty"`
	SetNumber int      `json:"setNumber"`
	WeightKg  *float64 `json:"weightKg,omitempty"` // 空文字→null→nil→層内で検証
	Reps      *int     `json:"reps,omitempty"`
	Note      *string  `json:"note,omitempty"`
}

// DomainToDTO converts domain.WorkoutRecord to WorkoutFormDTO
func WorkoutRecordToDTO(record *dom.WorkoutRecord) *WorkoutRecordDTO {
	if record == nil {
		return nil
	}

	dto := &WorkoutRecordDTO{
		ID:             domainIDToInt64Ptr(record.ID),
		PerformedDate:  record.PerformedDate.Format("2006-01-02"),
		StartedAt:      timeToHHmm(record.StartedAt),
		EndedAt:        timeToHHmm(record.EndedAt),
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
			isDefault := set.Exercise.IsPreset
			exerciseMap[exerciseID] = &ExerciseDTO{
				ID:            domainIDToInt64Ptr(&exerciseID),
				Name:          set.Exercise.Name,
				WorkoutPartID: domainIDToInt64Ptr(set.Exercise.PartID),
				IsDefault:     &isDefault,
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

func timeToHHmm(t *time.Time) *string {
	if t == nil {
		return nil
	}
	hhmm := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
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

	// Parse performedDate
	performedDate, err := time.Parse("2006-01-02", dto.PerformedDate)
	if err != nil {
		return nil, fmt.Errorf("invalid performedDate format: %w", err)
	}

	// Create WorkoutRecord
	record, err := dom.NewWorkoutRecord(dom.ULID(""), performedDate) // userID will be set by handler/usecase
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
		if exercise.IsDefault != nil {
			exerciseRef.IsPreset = *exercise.IsDefault
		}

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
func parseTimeWithDate(date time.Time, hhmmStr string) (time.Time, error) {
	t, err := time.Parse("15:04", hhmmStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(date.Year(), date.Month(), date.Day(), t.Hour(), t.Minute(), 0, 0, date.Location()), nil
}
