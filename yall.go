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

// Logger is the interface of a logger, that has methods to log with a context or without.
// With is used to decorate the logger with key-value pairs.
// ExecutionIDFrom allows to extract the execution id from the context, if present, and depends on the concrete logger
// implementation.
type Logger interface {
	ContextLogger
	NoContextLogger
	With(args ...interface{}) Logger
	ExecutionIDFrom(ctx context.Context) string
}

// ContextLogger is the interface of a logger, that has methods to log with a context.
// keysAndValues has to be formed by pairs of key and value, or single Fields.
type ContextLogger interface {
	Fatal(ctx context.Context, msg string, keysAndValues ...interface{})
	Panic(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, msg string, keysAndValues ...interface{})
	Debug(ctx context.Context, msg string, keysAndValues ...interface{})
}

// NoContextLogger is the interface of a logger, that has methods to log without a context.
// keysAndValues has to be formed by pairs of key and value, or single Fields.
type NoContextLogger interface {
	Fatalnc(msg string, keysAndValues ...interface{})
	Panicnc(msg string, keysAndValues ...interface{})
	Errornc(msg string, keysAndValues ...interface{})
	Warnnc(msg string, keysAndValues ...interface{})
	Infonc(msg string, keysAndValues ...interface{})
	Debugnc(msg string, keysAndValues ...interface{})
}

// Field is a loggable object, that can be used instead of key value pairs
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
