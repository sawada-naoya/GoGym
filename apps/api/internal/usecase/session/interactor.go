// 役割: セッションユースケースの実装（依存性注入とビジネスロジック）
// 受け取り: 外部サービスの実装（Repository, TokenService等）
// 処理: セッション管理のビジネスロジック実装
// 返却: UseCase インターフェースの実装
package session

import (
	"context"
	userUC "gogym-api/internal/usecase/user"
	"os"
	"time"

	"gogym-api/internal/adapter/dto"
	dom "gogym-api/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type interactor struct {
	// 外部依存関係
	ur UserRepository
	ph PasswordHasher
	// 他のユースケース
	uu userUC.UserUseCase
}

func NewInteractor(
	ur UserRepository,
	ph PasswordHasher,
	uu userUC.UserUseCase,
) SessionUseCase {
	return &interactor{
		ur: ur,
		ph: ph,
		uu: uu,
	}
}

func (i *interactor) Login(ctx context.Context, req dto.LoginRequest) error {

	// emailのバリデーション
	email, err := dom.NewEmail(req.Email)
	if err != nil {
		return err
	}

	// ユーザー検索
	user, err := i.ur.FindByEmail(ctx, email)
	if err != nil {
		return dom.NewDomainError("email_not_found")
	}
	if user == nil {
		return dom.NewDomainError("user_not_found")
	}

	// パスワード照合
	if err := i.ph.VerifyPassword(req.Password, user.PasswordHash); err != nil {
		return dom.NewDomainError("invalid_password")
	}
	return nil
}

func (i *interactor) CreateSession(ctx context.Context, email string) (dto.TokenResponse, error) {

	emailObj, err := dom.NewEmail(email)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	user, err := i.ur.FindByEmail(ctx, emailObj)
	if err != nil {
		return dto.TokenResponse{}, dom.NewDomainError("email_not_found")
	}
	if user == nil {
		return dto.TokenResponse{}, dom.NewDomainError("user_not_found")
	}

	now := time.Now()
	accessTTL := 24 * time.Hour

	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(accessTTL).Unix(),
		"iat": now.Unix(),
		"typ": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		AccessToken: tokenString,
		ExpiresIn:   int64(accessTTL.Seconds()),
	}, nil
}
