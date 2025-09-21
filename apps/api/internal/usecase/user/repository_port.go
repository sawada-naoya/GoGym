package user

import (
	"context"
	dom "gogym-api/internal/domain/user"
)

// Repository はユーザー関連のデータアクセスインターフェース
type Repository interface {
	Create(ctx context.Context, u *dom.User) error
	ExistsByEmail(ctx context.Context, email dom.Email) (bool, error)
}