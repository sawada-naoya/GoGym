package record

import (
	"time"

	"gorm.io/gorm"
)

// LocationRecord represents location table structure  
type LocationRecord struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	Address   string         `gorm:"not null"`
	Latitude  float64        `gorm:"not null"`
	Longitude float64        `gorm:"not null"`
	GymID     int64          `gorm:"not null;unique"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// TODO: Foreign key constraint
}

func (LocationRecord) TableName() string {
	return "locations"
}