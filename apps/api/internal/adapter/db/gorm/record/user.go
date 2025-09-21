package record

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"primaryKey;type:char(26)"` // ULID用
	Email        string         `gorm:"unique;not null;index"`
	PasswordHash string         `gorm:"not null"`
	Name         string         `gorm:"not null;column:name"` // DisplayName→Nameに統一
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
