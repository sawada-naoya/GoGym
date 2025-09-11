package record

import (
	"time"

	"gorm.io/gorm"
)

// RefreshTokenRecord represents refresh_tokens table structure
type RefreshTokenRecord struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	UserID    int64          `gorm:"not null;index"`
	TokenHash string         `gorm:"not null"`
	ExpiresAt time.Time      `gorm:"not null"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relations
	User *UserRecord `gorm:"foreignKey:UserID"`
}

func (RefreshTokenRecord) TableName() string {
	return "refresh_tokens"
}