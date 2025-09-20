package user

// PasswordHasher はパスワードをハッシュ化するインターフェース
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}

// PasswordService はパスワード関連サービスのインターフェース
type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}