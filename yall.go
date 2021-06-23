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
	contextLogger
	noContextLogger
	With(args ...interface{}) Logger
}

type contextLogger interface {
	Fatal(ctx context.Context, msg string, keysAndValues ...interface{})
	Panic(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
}

type noContextLogger interface {
	Fatalnc(msg string, keysAndValues ...interface{})
	Panicnc(msg string, keysAndValues ...interface{})
	Errornc(msg string, keysAndValues ...interface{})
	Warnnc(msg string, keysAndValues ...interface{})
	Infonc(msg string, keysAndValues ...interface{})
	Debugnc(msg string, keysAndValues ...interface{})
}

type Field struct {
	Name  string
	Value interface{}
}

// Error is a convenience method to make logging errors easier.
func Error(err error) *Field {
	return &Field{
		Name:  "error",
		Value: err,
	}
}
