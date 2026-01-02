package contextlib

import (
	"context"

	"github.com/sirupsen/logrus"
)

type userKey struct{}

var uk = userKey{}

func SetUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, uk, userID)
}

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(uk).(int64)
	return userID, ok
}

func MustGetUserID(ctx context.Context, l *logrus.Logger) int64 {
	userID, ok := GetUserID(ctx)
	if !ok {
		l.Fatal("failed to get user id from context")
	}
	return userID
}
