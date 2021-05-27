package zap_logger_test

import (
	"context"
	"testing"

	"github.com/phoops/yall"
	"github.com/phoops/yall/zap_logger"
	"github.com/stretchr/testify/assert"
)

func TestZapLoggerCreate(t *testing.T) {
	type testCase struct {
		name string
		opts []zap_logger.ZapLoggerOpt
	}

	testCases := []testCase{
		{
			name: "default",
			opts: []zap_logger.ZapLoggerOpt{},
		},
		{
			name: "production",
			opts: []zap_logger.ZapLoggerOpt{zap_logger.Production()},
		},
		{
			name: "name_key",
			opts: []zap_logger.ZapLoggerOpt{zap_logger.WithNameKey("test")},
		},
		{
			name: "execution_id_key",
			opts: []zap_logger.ZapLoggerOpt{zap_logger.WithExecutionIDKey("different_key")},
		},
		{
			name: "context_key",
			opts: []zap_logger.ZapLoggerOpt{zap_logger.WithExecutionIDContextKey("context_key")},
		},
		{
			name: "omit_missing_execution_id",
			opts: []zap_logger.ZapLoggerOpt{
				zap_logger.WithOmitExecutionIDWhenMissing(),
				zap_logger.WithExecutionIDContextKey("context_key"),
			},
		},
	}

	for _, tc := range testCases {
		logger, err := zap_logger.NewZapLogger("test", tc.opts...)
		var _ yall.Logger = logger

		assert.NoError(t, err)
		assert.NotNil(t, logger)

		logger = logger.With("test_name", tc.name)

		// should check the output, but at least thanks to zap DPANIC this would fail if we have a key without a value
		ctx := context.WithValue(context.Background(), yall.ExecutionIDKey, "request_id_test123")
		// logger.Fatal(ctx, "testing execution_id, fatal level")
		assert.Panics(t, func() {
			logger.Panic(ctx, "testing execution_id, panic level")
		})
		logger.Error(ctx, "testing execution_id, error level")
		logger.Warn(ctx, "testing execution_id, warn level")
		logger.Info(ctx, "testing execution_id, info level")
		logger.Debug(ctx, "testing execution_id, debug level")
	}
}
