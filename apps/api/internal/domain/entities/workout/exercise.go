package domain

// WorkoutExerciseRef represents a reference to a workout exercise
type WorkoutExerciseRef struct {
	ID     ID
	Name   string
	PartID *ID
	Owner  *ULID // nil ならプリセット、値があればユーザー作成
}
