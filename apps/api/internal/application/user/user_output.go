package user

import (
	"context"
	dom "gogym-api/internal/domain/entities/user"
)

// Repository はユーザーデータの永続化を担当
type Repository interface {
	Create(ctx context.Context, u *dom.User) error
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// PasswordHasher はパスワードのハッシュ化を担当
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}
