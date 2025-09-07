package middleware

import (
	"gogym-api/internal/adapter/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(e *echo.Echo, jwtService auth.TokenVerifier) {
	// ログを出力するミドルウェア
	e.Use(middleware.Logger())

	// panic時のリカバリミドルウェア
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// TODO: 本番では実際のフロントのURLを設定する
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	// JWT認証ミドルウェア（カスタム実装）
	e.Use(JWT(jwtService))
}
