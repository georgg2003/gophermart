package middleware

import (
	"net/http"
	"strings"

	"github.com/georgg2003/gophermart/internal/pkg/config"
	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/pkg/jwthelper"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func NewAuthMiddleware(
	cfg *config.Config,
	l *logrus.Logger,
	skipper func(c echo.Context) bool,
) echo.MiddlewareFunc {
	helper := jwthelper.New([]byte(cfg.JWTSecretKey))
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next(c)
			}
			req := c.Request()
			ctx := req.Context()
			var userID int64

			token := req.Header.Get(echo.HeaderAuthorization)
			if token == "" {
				l.WithContext(ctx).Error("authorization header is empty")
				c.Response().WriteHeader(http.StatusUnauthorized)
				return nil
			}
			if !strings.Contains(token, "Bearer ") {
				l.WithContext(ctx).Error("token is not bearer")
				c.String(http.StatusUnauthorized, "Invalid token")
				return nil
			}
			token = strings.Split(token, "Bearer ")[1]
			userID, err := helper.ReadAccessToken(token)
			if err != nil {
				l.WithContext(ctx).WithError(err).Error("got an invalid token")
				c.String(http.StatusUnauthorized, "Invalid token")
				return err
			}

			newCtx := contextlib.SetUserID(ctx, userID)
			c.SetRequest(req.WithContext(newCtx))
			return next(c)
		}
	}
}
