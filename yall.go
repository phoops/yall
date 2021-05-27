package yall

import (
	"context"
)

type executionIDType string

const (
	ExecutionIDKey executionIDType = "requestID"
)
const (
	MissingRequestIDKey = "missing_request_id"
)

type Logger interface {
	Fatal(ctx context.Context, msg string, keysAndValues ...interface{})
	Panic(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	With(args ...interface{}) Logger
}
