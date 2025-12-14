package session

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents refresh_tokens table structure（ULID対応）
type RefreshToken struct {
	JTI       string         `gorm:"primaryKey;type:char(26);column:jti"` // JWT ID (ULID)
	UserID    string         `gorm:"not null;index;type:char(26)"`        // User ID (ULID)
	RevokedAt *time.Time     `gorm:"index"`                               // 取り消し日時
	ExpiresAt time.Time      `gorm:"not null"`                            // 有効期限
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
