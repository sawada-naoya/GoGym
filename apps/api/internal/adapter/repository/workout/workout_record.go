package workout

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutRecord struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	UserID          string `gorm:"index"`
	GymID           *int64 `gorm:"index"` // gym_id追加
	PerformedDate   time.Time
	StartedAt       *time.Time
	EndedAt         *time.Time
	Note            *string
	ConditionLevel  *int
	DurationMinutes *int
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Record → Set（1:N）
	Sets []WorkoutSet `gorm:"foreignKey:WorkoutRecordID"`
	// Record → Gym（N:1）
	Gym *GymRecord `gorm:"foreignKey:GymID"`
}

type GymRecord struct {
	ID             int64  `gorm:"primaryKey"`
	Name           string `gorm:"size:255"`
	NormalizedName string `gorm:"size:255"`
}

func (GymRecord) TableName() string {
	return "gyms"
}

func (WorkoutRecord) TableName() string {
	return "workout_records"
}

type WorkoutSet struct {
	ID                int `gorm:"primaryKey;autoIncrement"`
	WorkoutRecordID   int `gorm:"index"`
	WorkoutExerciseID int `gorm:"index"`
	SetNumber         int
	WeightKg          float64
	Reps              int
	EstimatedMax      *float64
	Note              *string
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Set → Exercise（N:1）
	Exercise WorkoutExercise `gorm:"foreignKey:WorkoutExerciseID"`
}

func (WorkoutSet) TableName() string {
	return "workout_sets"
}

type WorkoutExercise struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	Name          string
	WorkoutPartID *int           `gorm:"index"`
	UserID        *string        `gorm:"index"` // nil ならプリセット
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Exercise → Part（N:1）
	Part *WorkoutPart `gorm:"foreignKey:WorkoutPartID"`
}

func (WorkoutExercise) TableName() string {
	return "workout_exercises"
}

type WorkoutPart struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Key       string
	UserID    *string        `gorm:"index"` // nil ならプリセット
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Exercises    []WorkoutExercise        `gorm:"foreignKey:WorkoutPartID"`
	Translations []WorkoutPartTranslation `gorm:"foreignKey:WorkoutPartID"`
}

func (WorkoutPart) TableName() string {
	return "workout_parts"
}

type WorkoutPartTranslation struct {
	ID            int `gorm:"primaryKey;autoIncrement"`
	WorkoutPartID int `gorm:"index"`
	Locale        string
	Name          string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

func (WorkoutPartTranslation) TableName() string {
	return "workout_part_translations"
}
