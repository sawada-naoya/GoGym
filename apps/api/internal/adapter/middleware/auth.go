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
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorization ヘッダーを取得
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "missing authorization header",
				})
			}

			// Bearer トークンを抽出
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "invalid authorization header format",
				})
			}

			tokenString := parts[1]
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "missing token",
				})
			}

			// トークンから user_id を抽出
			userID, err := extractUserIDFromToken(tokenString, jwtSecret)
			if err != nil {
				slog.Error("Authentication failed", "error", err.Error())
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			// Context に user_id を設定
			c.Set("user_id", userID)

			return next(c)
		}
	}
}

// extractUserIDFromToken は access_token から user_id を抽出します
func extractUserIDFromToken(tokenString string, jwtSecret string) (string, error) {
	// JWTトークンをパース（検証も同時に行う）
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	// クレームを取得
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// "sub" クレームからユーザーIDを取得
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", fmt.Errorf("sub claim not found in token")
	}

	return "", fmt.Errorf("invalid token")
}
