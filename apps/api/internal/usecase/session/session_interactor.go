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
	dom "gogym-api/internal/domain/entities"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
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
	accessTTL := 15 * time.Minute    // アクセストークンは短めに設定
	refreshTTL := 7 * 24 * time.Hour // リフレッシュトークンは7日間

	secret := []byte(os.Getenv("JWT_SECRET"))

	// アクセストークン生成
	accessClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(accessTTL).Unix(),
		"iat": now.Unix(),
		"typ": "access",
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// リフレッシュトークン生成
	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(refreshTTL).Unix(),
		"iat": now.Unix(),
		"typ": "refresh",
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email.String(),
		},
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}

func (i *interactor) RefreshToken(ctx context.Context, refreshToken string) (dto.TokenResponse, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	// リフレッシュトークンを検証
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, dom.NewDomainError("unexpected_signing_method")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return dto.TokenResponse{}, dom.NewDomainError("invalid_refresh_token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return dto.TokenResponse{}, dom.NewDomainError("invalid_token_claims")
	}

	// トークンタイプを確認
	if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
		return dto.TokenResponse{}, dom.NewDomainError("invalid_token_type")
	}

	// ユーザーIDを取得
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		return dto.TokenResponse{}, dom.NewDomainError("user_id_not_found")
	}

	// ULIDに変換
	userID, err := ulid.Parse(userIDStr)
	if err != nil {
		return dto.TokenResponse{}, dom.NewDomainError("invalid_user_id")
	}

	// ユーザー情報を取得
	user, err := i.ur.FindByID(ctx, userID)
	if err != nil || user == nil {
		return dto.TokenResponse{}, dom.NewDomainError("user_not_found")
	}

	// 新しいアクセストークンとリフレッシュトークンを生成
	now := time.Now()
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	// 新しいアクセストークン生成
	accessClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(accessTTL).Unix(),
		"iat": now.Unix(),
		"typ": "access",
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessTokenObj.SignedString(secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	// 新しいリフレッシュトークン生成
	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": now.Add(refreshTTL).Unix(),
		"iat": now.Unix(),
		"typ": "refresh",
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshTokenObj.SignedString(secret)
	if err != nil {
		return dto.TokenResponse{}, err
	}

	return dto.TokenResponse{
		User: dto.UserResponse{
			ID:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email.String(),
		},
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(accessTTL.Seconds()),
	}, nil
}
