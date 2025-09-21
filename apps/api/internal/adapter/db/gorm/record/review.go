package record

import (
	"time"

	"gorm.io/gorm"
)

// ReviewRecord represents review table structure
type ReviewRecord struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;index"`
	GymID     int64          `gorm:"not null;index"`
	Rating    int            `gorm:"not null"`
	Comment   *string        `gorm:"type:text"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ReviewRecord) TableName() string {
	return "reviews"
}
