package record

import (
	"time"

	"gorm.io/gorm"
)

// UserRecord represents user table structure
type UserRecord struct {
	ID           int64          `gorm:"primaryKey;autoIncrement"`
	Email        string         `gorm:"unique;not null;index"`
	PasswordHash string         `gorm:"not null"`
	DisplayName  string         `gorm:"not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (UserRecord) TableName() string {
	return "users"
}