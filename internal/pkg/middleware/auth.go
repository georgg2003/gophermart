package middleware

import (
	"net/http"
	"strings"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func NewAuthMiddleware(
	secretKey string,
	l *logrus.Logger,
	skipper func(c echo.Context) bool,
) echo.MiddlewareFunc {
	helper := jwthelper.New([]byte(secretKey))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper != nil && skipper(c) {
				return next(c)
			}
			req := c.Request()
			ctx := req.Context()
			var userID int64

			token := req.Header.Get(echo.HeaderAuthorization)
			if token == "" {
				l.WithContext(ctx).Error("authorization header is empty")
				return c.String(http.StatusUnauthorized, "No auth")
			}
			if !strings.Contains(token, "Bearer ") {
				l.WithContext(ctx).Error("token is not bearer")
				return c.String(http.StatusUnauthorized, "Invalid token")
			}
			token = strings.Split(token, "Bearer ")[1]
			userID, err := helper.ReadAccessToken(token)
			if err != nil {
				l.WithContext(ctx).WithError(err).Error("got an invalid token")
				return c.String(http.StatusUnauthorized, "Invalid token")
			}

			newCtx := contextlib.SetUserID(ctx, userID)
			c.SetRequest(req.WithContext(newCtx))
			return next(c)
		}
	}
}
