package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// LoggingMiddleware provides structured logging for requests
func LoggingMiddleware(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			start := time.Now()

			err := next(ctx)

			req := ctx.Request()
			res := ctx.Response()

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.String("remote_ip", ctx.RealIP()),
				zap.String("user_agent", req.UserAgent()),
				zap.Int("status", res.Status),
				zap.Int64("bytes_out", res.Size),
				zap.Duration("latency", time.Since(start)),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("Request failed", fields...)
			} else {
				logger.Info("Request completed", fields...)
			}

			return err
		}
	}
}
