// internal/adapter/http/server/provider.go
package server

import (
	"gogym-api/internal/configs"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho(httpCfg configs.HTTPConfig) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID()) // リクエストIDを付与
	e.Use(middleware.Secure())    // セキュリティヘッダを付与（XSS/Clickjacking などの軽減）
	e.Use(middleware.Gzip())      // Gzip 圧縮を有効化
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     httpCfg.CORS.AllowOrigins,
		AllowMethods:     httpCfg.CORS.AllowMethods,
		AllowHeaders:     httpCfg.CORS.AllowHeaders,
		AllowCredentials: httpCfg.CORS.AllowCreds,
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogLatency: true,
		LogMethod:  true,
		LogURI:     true,
		LogStatus:  true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			slog.Info("HTTP Request",
				"method", v.Method,
				"uri", v.URI,
				"status", v.Status,
				"latency", v.Latency.String(),
				"rid", c.Response().Header().Get(echo.HeaderXRequestID),
			)
			return nil
		},
	}))

	return e
}
