package middleware

import (
	"net/http"
	"strings"

	"gogym-api/internal/adapter/auth"

	"github.com/labstack/echo/v4"
)

// ctxKey はコンテキストキーの型定義
type ctxKey struct{}

// コンテキストに保存するキーの定数定義
const (
	CtxUserID = "uid"    // ユーザーIDを保存するキー
	CtxScopes = "scopes" // スコープ情報を保存するキー
)

// JWT はJWTトークン認証を行うミドルウェアを返す関数
// verifier: トークン検証を行うインターフェース
func JWT(verifier auth.TokenVerifier) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Authorizationヘッダーを取得
			h := c.Request().Header.Get("Authorization")
			
			// Bearer形式でない場合はエラー
			if !strings.HasPrefix(h, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing or invalid authorization header",
				})
			}
			
			// "Bearer "プレフィックスを削除してトークンを取得
			token := strings.TrimPrefix(h, "Bearer ")
			
			// トークンを検証してクレームを取得
			claims, err := verifier.Verify(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}
			
			// コンテキストにユーザー情報を保存
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxScopes, claims.Scopes)
			
			// 次のハンドラーを実行
			return next(c)
		}
	}
}
