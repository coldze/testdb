package utils

import (
	"context"
	"github.com/coldze/testdb/logs"
)

type loggerKey struct{}

type headerKey struct{}

var (
	//it is recommended to use structs as keys for values in context - not to overlap with other packages by accident.
	loggerCtxKey    loggerKey
	requestIdCtxKey headerKey

	//global variables are bad, but this one is not that bad - it's not exported outside and is used as a default logger, in case nothing was set in context - to remove checking == nil every single time.
	defaultLogger logs.Logger
)

func SetLogger(ctx context.Context, logger logs.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}

func GetLogger(ctx context.Context) logs.Logger {
	res := ctx.Value(loggerCtxKey)
	if res == nil {
		return defaultLogger
	}
	logger, ok := res.(logs.Logger)
	if !ok {
		return defaultLogger
	}
	return logger
}

func GetRequestID(ctx context.Context) string {
	res := ctx.Value(requestIdCtxKey)
	if res == nil {
		return ""
	}
	typedRes, ok := res.(string)
	if !ok {
		return ""
	}
	return typedRes
}

func SetRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIdCtxKey, id)
}

func init() {
	defaultLogger = logs.NewStdLogger()
}
