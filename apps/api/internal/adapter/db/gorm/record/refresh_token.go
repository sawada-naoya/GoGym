package record

import (
	"time"

	"gorm.io/gorm"
)

// RefreshTokenRecord represents refresh_tokens table structure（ULID対応）
type RefreshTokenRecord struct {
	ID        string         `gorm:"primaryKey;type:char(26)"` // ULID用
	UserID    string         `gorm:"not null;index;type:char(26)"` // User IDもULID
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