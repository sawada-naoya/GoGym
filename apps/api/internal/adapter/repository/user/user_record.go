package user

import (
	"time"

	"gorm.io/gorm"
)

// User はユーザーエンティティ用のGORMレコードを表す
type User struct {
	ID           string         `gorm:"primaryKey;type:char(26)"` // ULID用
	Email        string         `gorm:"unique;not null;index"`
	PasswordHash string         `gorm:"not null"`
	Name         string         `gorm:"not null;column:name"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
