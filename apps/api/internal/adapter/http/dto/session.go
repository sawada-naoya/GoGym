// 役割: セッション管理用のHTTP DTO（Data Transfer Object）
// 受け取り: HTTPリクエスト/レスポンスのJSON
// 処理: バリデーション、ドメインエンティティとの相互変換
// 返却: ドメイン層とHTTP層間のデータ変換結果
package dto

import (
	"time"
	"gogym-api/internal/domain/session"
)

// LoginResponse ログインレスポンス
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // 秒単位
}

// RefreshTokenRequest リフレッシュトークンリクエスト
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse リフレッシュトークンレスポンス
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // 秒単位
}

// ToRefreshTokenEntity パラメータからドメインエンティティに変換
func ToRefreshTokenEntity(jti, userID string, expiresAt, now time.Time) (*session.RefreshToken, error) {
	return session.NewRefreshToken(jti, userID, expiresAt, now)
}

// ToTokenPairValue パラメータからドメインバリューオブジェクトに変換
func ToTokenPairValue(accessToken, refreshToken string, expiresIn int) session.TokenPair {
	return session.NewTokenPair(accessToken, refreshToken, expiresIn)
}

// ToLoginResponse TokenPairからLoginResponse DTOに変換
func ToLoginResponse(tokenPair session.TokenPair) LoginResponse {
	return LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}
}

// ToRefreshTokenResponse TokenPairからRefreshTokenResponse DTOに変換
func ToRefreshTokenResponse(tokenPair session.TokenPair) RefreshTokenResponse {
	return RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    tokenPair.ExpiresIn,
	}
}
