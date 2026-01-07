package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			stop := time.Now()

			req := c.Request()
			res := c.Response()

			reqID := req.Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = res.Header().Get(echo.HeaderXRequestID)
			}

			entry := logger.WithFields(logrus.Fields{
				"remote_ip":  c.RealIP(),
				"host":       req.Host,
				"method":     req.Method,
				"path":       req.URL.Path,
				"status":     res.Status,
				"latency":    stop.Sub(start),
				"user_agent": req.UserAgent(),
				"request_id": reqID,
			})

			if err != nil {
				entry.WithField("error", err).Error("request completed with error")
				return err
			} else {
				entry.Info("request completed")
			}

			return nil
		}
	}
}
