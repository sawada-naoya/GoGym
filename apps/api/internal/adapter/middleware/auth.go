package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
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

		token := parts[1]
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
				"error": "missing token",
			})
		}

		// トークンから user_id を抽出
		// 注意: これは簡易実装です
		// 本番環境では、バックエンド API から返却された access_token を使用し、
		// そのトークンを検証する必要があります
		userID := extractUserIDFromToken(token)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
				"error": "invalid token",
			})
		}

		// Context に user_id を設定
		c.Set("user_id", userID)

		return next(c)
	}
}

// extractUserIDFromToken は access_token から user_id を抽出します
//
// 現在の実装:
// - NextAuth の JWT には直接 user_id が含まれていない
// - バックエンドログイン時に返却された access_token を使用する想定
// - この access_token を検証・デコードして user_id を取得する
//
// TODO: 実際の JWT 検証処理を実装
// - JWT ライブラリ (github.com/golang-jwt/jwt/v5) を使用
// - トークンの署名を検証
// - 有効期限をチェック
// - user_id を抽出
func extractUserIDFromToken(token string) string {
	// 暫定実装：開発用
	// 実際には JWT をデコード・検証する必要があります

	// 仮実装: トークンをそのまま user_id として返す
	// これは開発・テスト用の簡易実装です
	return token
}
