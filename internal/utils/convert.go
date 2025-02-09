package utils

import (
	"context"
	"time"
)

type ctxKey string

var (
	UserIdCtxKey      = ctxKey("userID")
	DataLoadersCtxKey = ctxKey("dataloaders")
)

func ConvertTimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func UserIDFromContext(ctx context.Context) string {
	return ctx.Value(UserIdCtxKey).(string)
}
