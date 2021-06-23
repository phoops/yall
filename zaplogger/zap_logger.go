// Package zap_logger provides an implementation of the Yall logger interface
package zaplogger

import (
	"context"

	"github.com/phoops/yall"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is just a thin wrapper around a sugared zap logger with some opinionated defaults.
type Logger struct {
	l    zap.SugaredLogger
	conf *loggerConf
}

// NewLogger creates a new Logger, configured with the passed LoggerOpt args
// By default, with no options, the logger is configured for development.
func NewLogger(name string, opts ...LoggerOpt) (yall.Logger, error) {
	conf := defaultConf()

	for _, opt := range opts {
		conf = opt(conf)
	}

	var logger *zap.Logger
	var err error
	// Check for environment
	if conf.production {
		conf := zap.NewProductionConfig()
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = conf.Build()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, errors.Wrap(err, "error during log initialization")
	}
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	sugared := logger.Sugar().With(conf.nameKey, name)
	return &Logger{
		l:    *sugared,
		conf: conf,
	}, nil
}

type LoggerOpt func(opts *loggerConf) *loggerConf

// Production uses a zap production config, with ISO 8601 timestamps.
func Production() LoggerOpt {
	return func(opts *loggerConf) *loggerConf {
		opts.production = true

		return opts
	}
}

// WithNameKey configures the key to use for the logger's name.
func WithNameKey(nameKey string) LoggerOpt {
	return func(opts *loggerConf) *loggerConf {
		opts.nameKey = nameKey

		return opts
	}
}

// WithExecutionIDKey sets the key to use to log the execution id.
func WithExecutionIDKey(requestIDKey string) LoggerOpt {
	return func(opts *loggerConf) *loggerConf {
		opts.executionIDKey = requestIDKey

		return opts
	}
}

// WithExecutionIDContextKey configures the key to use to extract execution id from context.
func WithExecutionIDContextKey(requestIDContextKey interface{}) LoggerOpt {
	return func(opts *loggerConf) *loggerConf {
		opts.executionIDContextKey = requestIDContextKey

		return opts
	}
}

func WithOmitExecutionIDWhenMissing() LoggerOpt {
	return func(opts *loggerConf) *loggerConf {
		opts.omitExecutionIDWhenMissing = true

		return opts
	}
}

type loggerConf struct {
	production                 bool
	nameKey                    string
	executionIDKey             string
	executionIDContextKey      interface{}
	omitExecutionIDWhenMissing bool
}

func defaultConf() *loggerConf {
	return &loggerConf{
		production:                 false,
		nameKey:                    "service",
		executionIDKey:             string(yall.ExecutionIDKey),
		executionIDContextKey:      yall.ExecutionIDKey,
		omitExecutionIDWhenMissing: false,
	}
}

func (l *Logger) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Panic(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Panicw(msg, keysAndValues...)
}

func (l *Logger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Errorw(msg, keysAndValues...)
}

func (l *Logger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Warnw(msg, keysAndValues...)
}

func (l *Logger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Infow(msg, keysAndValues...)
}

func (l *Logger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	keysAndValues = extractFields(keysAndValues...)
	l.l.Debugw(msg, keysAndValues...)
}

func (l *Logger) Fatalnc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Fatalw(msg, keysAndValues...)
}

func (l *Logger) Panicnc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Panicw(msg, keysAndValues...)
}

func (l *Logger) Errornc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Errorw(msg, keysAndValues...)
}

func (l *Logger) Warnnc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Warnw(msg, keysAndValues...)
}

func (l *Logger) Infonc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Infow(msg, keysAndValues...)
}

func (l *Logger) Debugnc(msg string, keysAndValues ...interface{}) {
	keysAndValues = extractFields(keysAndValues...)
	l.l.Debugw(msg, keysAndValues...)
}

func (l *Logger) With(args ...interface{}) yall.Logger {
	return &Logger{
		l:    *l.l.With(args...),
		conf: l.conf,
	}
}

func (l *Logger) addExecutionIDField(ctx context.Context, keysAndValues ...interface{}) []interface{} {
	executionID := l.ExecutionIDFrom(ctx)
	if executionID != yall.MissingExecutionID || !l.conf.omitExecutionIDWhenMissing {
		keysAndValues = append([]interface{}{zap.String(l.conf.executionIDKey, executionID)}, keysAndValues...)
	}
	return keysAndValues
}

// ExecutionIDFrom retrieves the execution id from the given context, if present.
func (l *Logger) ExecutionIDFrom(ctx context.Context) string {
	if ctx == nil {
		return yall.MissingExecutionID
	}
	reqID := ctx.Value(l.conf.executionIDContextKey)
	if reqID == nil {
		return yall.MissingExecutionID
	}
	if reqID, ok := reqID.(string); ok {
		return reqID
	}

	return yall.MissingExecutionID
}

func extractFields(keysAndValues ...interface{}) []interface{} {
	var processedKeysAndValues []interface{}

	for _, keyOrValue := range keysAndValues {
		if field, ok := keyOrValue.(*yall.Field); ok {
			processedKeysAndValues = append(processedKeysAndValues, field.Name, field.Value)
		} else {
			processedKeysAndValues = append(processedKeysAndValues, keyOrValue)
		}
	}

	return processedKeysAndValues
}
