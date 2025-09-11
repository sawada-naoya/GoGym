package record

import (
	"time"

	"gorm.io/gorm"
)

// UserRecord represents user table structure
type UserRecord struct {
	ID              int64          `gorm:"primaryKey;autoIncrement"`
	Email           string         `gorm:"unique;not null;index"`
	CryptedPassword string         `gorm:"not null"`
	Salt            string         `gorm:"not null"`
	Name            string         `gorm:"not null"`
	CreatedAt       time.Time      `gorm:"not null"`
	UpdatedAt       time.Time      `gorm:"not null"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// TODO: Relations
}

func (UserRecord) TableName() string {
	return "users"
}