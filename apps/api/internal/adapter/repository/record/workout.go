package record

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutRecord struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	UserID          string `gorm:"index"`
	PerformedDate   time.Time
	StartedAt       *time.Time
	EndedAt         *time.Time
	Place           *string
	Note            *string
	ConditionLevel  *int
	DurationMinutes *int
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Record → Set（1:N）
	Sets []WorkoutSet `gorm:"foreignKey:WorkoutRecordID"`
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

type WorkoutPart struct {
	ID        int `gorm:"primaryKey;autoIncrement"`
	Name      string
	UserID    *string        `gorm:"index"` // nil ならプリセット
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Exercises []WorkoutExercise `gorm:"foreignKey:WorkoutPartID"`
}
