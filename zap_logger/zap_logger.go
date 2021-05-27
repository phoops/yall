package zap_logger

import (
	"context"

	"github.com/phoops/yall"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

func Error(err error) Field {
	return zap.Error(err)
}

// ZapLogger is just a thin wrapper around a sugared zap logger with some opinionated defaults
type ZapLogger struct {
	l    zap.SugaredLogger
	conf *zapLoggerConf
}

// NewZapLogger creates a new ZapLogger, configured with the passed ZapLoggerOpt args
// By default, with no options, the logger is configured for development.
func NewZapLogger(name string, opts ...ZapLoggerOpt) (yall.Logger, error) {
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
	return &ZapLogger{
		l:    *sugared,
		conf: conf,
	}, nil
}

type ZapLoggerOpt func(opts *zapLoggerConf) *zapLoggerConf

// Production uses a zap production config, with ISO 8601 timestamps
func Production() ZapLoggerOpt {
	return func(opts *zapLoggerConf) *zapLoggerConf {
		opts.production = true
		return opts
	}
}

// WithNameKey configures the key to use for the logger's name
func WithNameKey(nameKey string) ZapLoggerOpt {
	return func(opts *zapLoggerConf) *zapLoggerConf {
		opts.nameKey = nameKey
		return opts
	}
}

// WithExecutionIDKey sets the key to use to log the request id
func WithExecutionIDKey(requestIDKey string) ZapLoggerOpt {
	return func(opts *zapLoggerConf) *zapLoggerConf {
		opts.executionIDKey = requestIDKey
		return opts
	}
}

// WithExecutionIDContextKey configures the key to use to extract request id from context
func WithExecutionIDContextKey(requestIDContextKey interface{}) ZapLoggerOpt {
	return func(opts *zapLoggerConf) *zapLoggerConf {
		opts.executionIDContextKey = requestIDContextKey
		return opts
	}
}

func WithOmitExecutionIDWhenMissing() ZapLoggerOpt {
	return func(opts *zapLoggerConf) *zapLoggerConf {
		opts.omitExecutionIDWhenMissing = true
		return opts
	}
}

type zapLoggerConf struct {
	production                 bool
	nameKey                    string
	executionIDKey             string
	executionIDContextKey      interface{}
	omitExecutionIDWhenMissing bool
}

func defaultConf() *zapLoggerConf {
	return &zapLoggerConf{
		production:                 false,
		nameKey:                    "service",
		executionIDKey:             string(yall.ExecutionIDKey),
		executionIDContextKey:      yall.ExecutionIDKey,
		omitExecutionIDWhenMissing: false,
	}
}

func (l *ZapLogger) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	executionID := l.getExecutionIDFromContext(ctx)
	if executionID != yall.MissingRequestIDKey || !l.conf.omitExecutionIDWhenMissing {
		keysAndValues = append([]interface{}{zap.String(l.conf.executionIDKey, executionID)}, keysAndValues...)
	}

	l.l.Fatalw(msg, keysAndValues...)
}

func (l *ZapLogger) Panic(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	l.l.Panicw(msg, keysAndValues...)
}

func (l *ZapLogger) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	l.l.Errorw(msg, keysAndValues...)
}

func (l *ZapLogger) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	l.l.Warnw(msg, keysAndValues...)
}

func (l *ZapLogger) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	l.l.Infow(msg, keysAndValues...)
}

func (l *ZapLogger) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	keysAndValues = l.addExecutionIDField(ctx, keysAndValues...)
	l.l.Debugw(msg, keysAndValues...)
}

func (l *ZapLogger) With(args ...interface{}) yall.Logger {
	return &ZapLogger{
		l:    *l.l.With(args...),
		conf: l.conf,
	}
}

func (l *ZapLogger) addExecutionIDField(ctx context.Context, keysAndValues ...interface{}) []interface{} {
	executionID := l.getExecutionIDFromContext(ctx)
	if executionID != yall.MissingRequestIDKey || !l.conf.omitExecutionIDWhenMissing {
		keysAndValues = append([]interface{}{zap.String(l.conf.executionIDKey, executionID)}, keysAndValues...)
	}
	return keysAndValues
}

func (l *ZapLogger) getExecutionIDFromContext(ctx context.Context) string {
	reqID := ctx.Value(l.conf.executionIDContextKey)
	if reqID == nil {
		return yall.MissingRequestIDKey
	}
	if reqID, ok := reqID.(string); ok {
		return reqID
	} else {
		return yall.MissingRequestIDKey
	}
}
