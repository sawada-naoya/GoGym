package domain

// WorkoutPart represents a workout body part (e.g., chest, back, legs)
type WorkoutPart struct {
	ID           ID
	Key          string
	Owner        *ULID                    // nil ならプリセット、値があればユーザー作成
	Translations []WorkoutPartTranslation // 多言語対応
	Exercises    []WorkoutExerciseRef     // この部位に紐づく種目
}

// WorkoutPartTranslation represents a translation of a workout part name
type WorkoutPartTranslation struct {
	ID            ID
	WorkoutPartID ID
	Locale        string
	Name          string
}
