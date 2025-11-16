package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware はJWT認証ミドルウェアを生成します
func AuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	slog.Info("[AuthMiddleware] Initializing with JWT secret", "secret_length", len(jwtSecret))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()

			// Authorization ヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				slog.ErrorContext(ctx, "[AuthMiddleware] Missing Authorization header")
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			// Bearer トークンを抽出
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				slog.ErrorContext(ctx, "[AuthMiddleware] Invalid Authorization header format", "header", authHeader)
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]
			if tokenString == "" {
				slog.ErrorContext(ctx, "[AuthMiddleware] Missing token")
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "missing token",
				})
			}

			slog.InfoContext(ctx, "[AuthMiddleware] Token received", "token_prefix", tokenString[:20]+"...")

			// トークンから user_id を抽出
			userID, err := extractUserIDFromToken(tokenString, jwtSecret)
			if err != nil {
				slog.ErrorContext(ctx, "[AuthMiddleware] Failed to extract user ID from token", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": fmt.Sprintf("invalid token: %v", err),
				})
			}

			slog.InfoContext(ctx, "[AuthMiddleware] Successfully authenticated", "userID", userID)

			// Context に user_id を設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// extractUserIDFromToken は access_token から user_id を抽出します
func extractUserIDFromToken(tokenString string, jwtSecret string) (string, error) {
	// JWTトークンをパース
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			slog.Error("[extractUserIDFromToken] Unexpected signing method", "method", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		slog.Info("[extractUserIDFromToken] JWT signing method verified", "method", token.Method.Alg())
		return []byte(jwtSecret), nil
	})

	if err != nil {
		slog.Error("[extractUserIDFromToken] Failed to parse token", "error", err)
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	slog.Info("[extractUserIDFromToken] Token parsed successfully", "valid", token.Valid)

	// クレームを取得
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		slog.Info("[extractUserIDFromToken] Token claims retrieved", "claims_count", len(claims))

		// "sub" クレームからユーザーIDを取得
		if sub, ok := claims["sub"].(string); ok {
			slog.Info("[extractUserIDFromToken] User ID extracted from sub claim", "userID", sub)
			return sub, nil
		}
		slog.Error("[extractUserIDFromToken] Sub claim not found or not a string", "claims", claims)
		return "", fmt.Errorf("sub claim not found in token")
	}

	slog.Error("[extractUserIDFromToken] Invalid token claims or token not valid")
	return "", fmt.Errorf("invalid token claims")
}
