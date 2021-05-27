package yall

import (
	"context"
)

type executionIDType string

const (
	ExecutionIDKey executionIDType = "executionID"
)
const (
	MissingExecutionID = "missing_execution_id"
)

type Logger interface {
	Fatal(ctx context.Context, msg string, keysAndValues ...interface{})
	Fatalnc(msg string, keysAndValues ...interface{})
	Panic(ctx context.Context, msg string, keysAndValues ...interface{})
	Panicnc(msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Errornc(msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Warnnc(msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Infonc(msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
	Debugnc(msg string, keysAndValues ...interface{})
	With(args ...interface{}) Logger
}