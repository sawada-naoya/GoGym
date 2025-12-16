package gym

import (
	"time"

	"gorm.io/gorm"
)

type GymRecord struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	Name            string         `gorm:"size:255;not null;index:idx_gyms_name"`
	Latitude        float64        `gorm:"type:decimal(10,7);not null;index:idx_gyms_location,priority:1"`
	Longitude       float64        `gorm:"type:decimal(10,7);not null;index:idx_gyms_location,priority:2"`
	SourceURL       string         `gorm:"size:1000;not null"`
	PrimaryPhotoURL *string        `gorm:"size:1000"`
	PlaceID         *string        `gorm:"size:128;uniqueIndex:uq_gyms_place_id"`
	CreatedBy       string         `gorm:"size:26;not null"`
	UpdatedBy       *string        `gorm:"size:26"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (GymRecord) TableName() string {
	return "gyms"
}
