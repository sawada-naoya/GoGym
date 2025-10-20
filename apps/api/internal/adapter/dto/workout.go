package dto

import (
	"fmt"
	"time"

	dom "gogym-api/internal/domain/workout"
)

type WorkoutFormDTO struct {
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
func WorkoutRecordToDTO(record *dom.WorkoutRecord) *WorkoutFormDTO {
	if record == nil {
		return nil
	}

	dto := &WorkoutFormDTO{
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
