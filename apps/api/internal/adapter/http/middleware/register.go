package middleware

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupMiddleware(e *echo.Echo) {
	// ログを出力するミドルウェア
	e.Use(middleware.Logger())

	// panic時のリカバリミドルウェア
	e.Use(middleware.Recover())

	frontOrigin := os.Getenv("FRONT_ORIGIN")
	if frontOrigin == "" {
		frontOrigin = "http://localhost:3003"
	}

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontOrigin},
		AllowMethods: []string{
			http.MethodGet, http.MethodHead, http.MethodPost,
			http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
		},
		AllowCredentials: true,
	}))

}
