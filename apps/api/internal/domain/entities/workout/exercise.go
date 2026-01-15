package workout

import (
	dom "gogym-api/internal/domain/entities"
)

// WorkoutExerciseRef represents a reference to a workout exercise
type WorkoutExerciseRef struct {
	ID     dom.ID
	Name   string
	PartID *dom.ID
	Owner  *dom.ULID // nil ならプリセット、値があればユーザー作成
}
