package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"gogym-api/configs"
	"gogym-api/internal/domain/user"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims represents JWT claims
type Claims struct {
	UserID user.ID `json:"user_id"`
	Email  string    `json:"email"`
	Type   string    `json:"type"` // "access" or "refresh"
	JTI    string    `json:"jti"`  // JWT ID for refresh token tracking
	jwt.RegisteredClaims
}

// TokenService implements JWT token operations
type TokenService struct {
	accessSecret     []byte
	refreshSecret    []byte
	accessExpiresIn  time.Duration
	refreshExpiresIn time.Duration
	issuer           string
}

// NewTokenService creates a new token service from auth config
func NewTokenService(authConfig configs.AuthConfig) *TokenService {
	return &TokenService{
		accessSecret:     []byte(authConfig.JWTSecret),
		refreshSecret:    []byte(authConfig.JWTSecret), // 同じシークレットを使用
		accessExpiresIn:  authConfig.AccessExpiresIn,
		refreshExpiresIn: authConfig.RefreshExpiresIn,
		issuer:           authConfig.Issuer,
	}
}

// GenerateTokens generates access and refresh tokens
func (s *TokenService) GenerateTokens(userID user.ID, email string) (string, string, error) {
	now := time.Now()
	jti := uuid.New().String()

	// Generate access token
	accessClaims := Claims{
		UserID: userID,
		Email:  email,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("user:%d", userID),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.accessSecret)
	if err != nil {
		return "", "", user.NewDomainError(user.ErrInternal, "token_signing_failed", "failed to sign access token")
	}

	// Generate refresh token
	refreshClaims := Claims{
		UserID: userID,
		Email:  email,
		Type:   "refresh",
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("user:%d", userID),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.refreshSecret)
	if err != nil {
		return "", "", user.NewDomainError(user.ErrInternal, "token_signing_failed", "failed to sign refresh token")
	}

	return accessTokenString, refreshTokenString, nil
}

// ValidateAccessToken validates and parses access token
func (s *TokenService) ValidateAccessToken(tokenString string) (user.ID, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.accessSecret, nil
	})

	if err != nil {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "invalid_access_token", "invalid access token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "invalid_access_token", "invalid access token claims")
	}

	if claims.Type != "access" {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "wrong_token_type", "wrong token type")
	}

	return claims.UserID, claims.Email, nil
}

// ValidateRefreshToken validates and parses refresh token
func (s *TokenService) ValidateRefreshToken(tokenString string) (user.ID, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.refreshSecret, nil
	})

	if err != nil {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "invalid_refresh_token", "invalid refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "invalid_refresh_token", "invalid refresh token claims")
	}

	if claims.Type != "refresh" {
		return 0, "", user.NewDomainError(user.ErrUnauthorized, "wrong_token_type", "wrong token type")
	}

	// Create token hash from JTI for database storage/lookup
	tokenHash := s.createTokenHash(claims.JTI)

	return claims.UserID, tokenHash, nil
}

// createTokenHash creates a hash from JWT ID for database storage
func (s *TokenService) createTokenHash(jti string) string {
	hash := sha256.Sum256([]byte(jti))
	return hex.EncodeToString(hash[:])
}

// ExtractUserFromToken extracts user information from access token without full validation
// This is useful for logging/tracing purposes
func (s *TokenService) ExtractUserFromToken(tokenString string) (user.ID, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.accessSecret, nil
	})

	if err != nil {
		return 0, "", err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, claims.Email, nil
}
