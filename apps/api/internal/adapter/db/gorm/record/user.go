// internal/adapter/db/gorm/record/user.go
// 役割: ユーザー関連のGORMレコード構造体（Infrastructure Layer）
// DB行の形でGORMタグ付きstruct。ドメインエンティティとの変換はconverterで実行
package record

import (
	"time"
)

// UserRecord はユーザーエンティティ用のGORMレコードを表す
type UserRecord struct {
	ID           int64     `gorm:"primaryKey;autoIncrement"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	DisplayName  string    `gorm:"size:100;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName はGORM用のテーブル名を返す
func (UserRecord) TableName() string {
	return "users"
}

// RefreshTokenRecord はリフレッシュトークンエンティティ用のGORMレコードを表す
type RefreshTokenRecord struct {
	ID        int64       `gorm:"primaryKey;autoIncrement"`
	UserID    int64       `gorm:"not null;index"`
	TokenHash string      `gorm:"uniqueIndex;size:255;not null"`
	ExpiresAt time.Time   `gorm:"not null;index"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
	User      *UserRecord `gorm:"foreignKey:UserID"`
}

// TableName はGORM用のテーブル名を返す
func (RefreshTokenRecord) TableName() string {
	return "refresh_tokens"
}
