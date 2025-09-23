package user

import (
	"context"
	dom "gogym-api/internal/domain/user"
)

// usecase → output (repository, external services, etc.)

// Repository はユーザーデータの永続化を担当
type Repository interface {
	Create(ctx context.Context, u *dom.User) error
	FindByEmail(ctx context.Context, email dom.Email) (*dom.User, error)
	ExistsByEmail(ctx context.Context, email dom.Email) (bool, error)
}

// PasswordHasher はパスワードのハッシュ化を担当
type PasswordHasher interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) error
}

// IDProvider はユニークIDの生成を担当
type IDProvider interface {
	NewUserID() string
}

// TokenService はJWTトークンの生成・検証を担当
type TokenService interface {
	GenerateTokens(userID dom.ID, email string) (string, string, error)
	ValidateAccessToken(tokenString string) (dom.ID, string, error)
	ValidateRefreshToken(tokenString string) (dom.ID, string, error)
}