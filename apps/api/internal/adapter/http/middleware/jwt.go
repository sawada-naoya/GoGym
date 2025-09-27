package middleware

import (
	"net/http"
	"strings"

	jwt "gogym-api/internal/adapter/auth"

	"github.com/labstack/echo/v4"
)

const (
	CtxUserID = "uid"
)

type AccessParser interface {
	ParseRefresh(tokenStr string) (jwt.RefreshClaims, error)
}

func JWT(verifier AccessParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return c.NoContent(http.StatusUnauthorized)
			}
			token := strings.TrimPrefix(h, "Bearer ")
			if token == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			claims, err := verifier.ParseRefresh(token)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set(CtxUserID, claims.UserID)

			return next(c)
		}
	}
}
