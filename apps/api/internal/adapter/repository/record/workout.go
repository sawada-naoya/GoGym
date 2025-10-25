package record

import (
	"time"

	"gorm.io/gorm"
)

type WorkoutRecord struct {
	ID              int            `gorm:"primaryKey;autoIncrement"`
	UserID          string         `gorm:"type:char(26);not null;index"`
	PerformedDate   time.Time      `gorm:"type:date;not null"`
	StartedAt       *time.Time     `gorm:"type:datetime"`
	EndedAt         *time.Time     `gorm:"type:datetime"`
	Place           *string        `gorm:"type:varchar(100)"`
	Note            *string        `gorm:"type:text"`
	ConditionLevel  *int           `gorm:"type:tinyint"`
	DurationMinutes *int           `gorm:"type:int"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Record → Set（1:N）
	Sets []WorkoutSet `gorm:"foreignKey:WorkoutRecordID;constraint:OnDelete:CASCADE"`
}

type WorkoutSet struct {
	ID                int            `gorm:"primaryKey;autoIncrement"`
	WorkoutRecordID   int            `gorm:"not null;index"`
	WorkoutExerciseID int            `gorm:"not null;index"`
	SetNumber         int            `gorm:"not null"`
	WeightKg          float64        `gorm:"type:decimal(6,2);not null;default:0.00"`
	Reps              int            `gorm:"not null;default:0"`
	EstimatedMax      *float64       `gorm:"type:decimal(6,2)"`
	Note              *string        `gorm:"type:varchar(255)"`
	CreatedAt         time.Time      `gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Set → Exercise（N:1）
	Exercise WorkoutExercise `gorm:"foreignKey:WorkoutExerciseID"`
}

type WorkoutExercise struct {
	ID            int            `gorm:"primaryKey;autoIncrement"`
	Name          string         `gorm:"type:varchar(100);not null"`
	WorkoutPartID *int           `gorm:"index"`
	IsDefault     bool           `gorm:"type:tinyint(1);not null;default:0"`
	UserID        *string        `gorm:"type:char(26);index"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`

	// Exercise → Part（N:1）
	Part *WorkoutPart `gorm:"foreignKey:WorkoutPartID"`
}

type WorkoutPart struct {
	ID        int            `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"type:varchar(50);not null"`
	IsDefault bool           `gorm:"type:tinyint(1);not null;default:0"`
	UserID    *string        `gorm:"type:char(26);index"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Exercises []WorkoutExercise `gorm:"foreignKey:WorkoutPartID"`
}
