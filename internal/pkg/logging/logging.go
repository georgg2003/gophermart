package logging

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/georgg2003/gophermart/internal/pkg/contextlib"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) clone() *Logger {
	c := *l
	return &c
}

func (l *Logger) WithRequestCtx(ctx context.Context) *Logger {
	reqInfo, ok := contextlib.GetRequestInfo(ctx)
	if !ok {
		return l
	}
	userIDStr := ""
	userID, ok := contextlib.GetUserID(ctx)
	if ok {
		userIDStr = fmt.Sprint(userID)
	}
	newLogger := l.clone()
	newLogger.Logger = newLogger.Logger.With(
		slog.String("request_id", reqInfo.RequestID),
		slog.String("remote_ip", reqInfo.RemoteIP),
		slog.String("method", reqInfo.Method),
		slog.String("path", reqInfo.Path),
		slog.String("user_agent", reqInfo.UserAgent),
		slog.String("user_id", userIDStr),
	)
	return newLogger
}

func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	newLogger := l.clone()
	newLogger.Logger = newLogger.Logger.With(
		slog.String("error", err.Error()),
	)
	return newLogger
}

func (l *Logger) With(attrs ...any) *Logger {
	if len(attrs) == 0 {
		return l
	}
	newLogger := l.clone()
	newLogger.Logger = newLogger.Logger.With(attrs...)
	return newLogger
}

func (l *Logger) WithString(key string, value string) *Logger {
	if key == "" {
		return l
	}
	newLogger := l.clone()
	newLogger.Logger = newLogger.Logger.With(slog.String(key, value))
	return newLogger
}

func (l *Logger) WithGroup(name string) *Logger {
	if name == "" {
		return l
	}
	newLogger := l.clone()
	newLogger.Logger = newLogger.Logger.WithGroup(name)
	return newLogger
}

func (l *Logger) Fatal(msg string, attrs ...any) {
	l.Error(msg, attrs...)
	os.Exit(1)
}

func (l *Logger) WithLevel() {
	l = l.clone()
	l.Logger.Handler()
}

type LoggerOption func(*Logger)

func New(level slog.Level) *Logger {
	return &Logger{
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})),
	}
}
