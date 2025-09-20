package user

import dom "gogym-api/internal/domain/user"

// TokenService はトークン関連サービスのインターフェース
type TokenService interface {
	GenerateTokens(userID dom.ID, email string) (string, string, error)
	ValidateAccessToken(tokenString string) (dom.ID, string, error)
	ValidateRefreshToken(tokenString string) (dom.ID, string, error)
}