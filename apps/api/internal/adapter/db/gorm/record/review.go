package record

import (
	"time"

	"gorm.io/gorm"
)

// ReviewRecord represents review table structure
type ReviewRecord struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Title     string         `gorm:"not null"`
	Content   string         `gorm:"type:text"`
	Rating    int            `gorm:"not null"`
	ImageURL  *string
	GymID     int64          `gorm:"not null;index"`
	UserID    int64          `gorm:"not null;index"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// TODO: Foreign key relations
}

func (ReviewRecord) TableName() string {
	return "reviews"
}