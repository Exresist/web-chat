package ctxkey

import (
	"context"

	log "go.uber.org/zap"
)

type ctxKey int

const (
	ctxKeyRequestID ctxKey = iota
	ctxKeyLogger
)

func GetRequestID(ctx context.Context) string {
	return ctx.Value(ctxKeyRequestID).(string)
}

func PutRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, reqID)
}

func GetLogger(ctx context.Context) *log.SugaredLogger {
	return ctx.Value(ctxKeyLogger).(*log.SugaredLogger)
}

func PutLogger(ctx context.Context, logger *log.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxKeyLogger, logger)
}
