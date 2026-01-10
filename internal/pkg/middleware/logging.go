package middleware

import (
	"fmt"
	"time"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
	"github.com/georgg2003/gophermart/internal/pkg/logging"
	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(logger *logging.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			reqID := req.Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = res.Header().Get(echo.HeaderXRequestID)
			}

			ctx := req.Context()
			ctx = contextlib.SetRequestInfo(ctx, contextlib.RequestInfo{
				RequestID: reqID,
				RemoteIP:  c.RealIP(),
				Host:      req.Host,
				Method:    req.Method,
				Path:      req.URL.Path,
				UserAgent: req.UserAgent(),
			})
			c.SetRequest(req.WithContext(ctx))

			err := next(c)
			stop := time.Now()

			logger = logger.WithRequestCtx(ctx).WithString(
				"latency", fmt.Sprint(stop.Sub(start)),
			)

			if err != nil {
				logger.WithError(err).Error("request completed with error")
				return err
			} else {
				logger.Info("request completed")
			}

			return nil
		}
	}
}
