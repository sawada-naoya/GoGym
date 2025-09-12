package entity

import "time"

// BaseEntity は全エンティティの共通フィールド
type BaseEntity struct {
	ID        ID
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ID は共通のIDタイプ
type ID int64

// User はユーザードメインエンティティ
type User struct {
	BaseEntity
	Email        Email
	PasswordHash string
	DisplayName  string
}

// Email はメールアドレスを表すバリューオブジェクト
type Email struct {
	value string
}

// NewEmail はEmailを作成する
func NewEmail(email string) (Email, error) {
	// TODO: バリデーションを実装
	return Email{value: email}, nil
}

// String はメールアドレスの文字列表現を返す
func (e Email) String() string {
	return e.value
}

// RefreshToken はリフレッシュトークンエンティティ
type RefreshToken struct {
	ID        ID
	UserID    ID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
	
	// Relations
	User *User
}