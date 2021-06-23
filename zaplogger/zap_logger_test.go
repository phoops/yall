package zaplogger_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/phoops/yall"
	"github.com/phoops/yall/zaplogger"
	"github.com/stretchr/testify/assert"
)

func TestZapLoggerCreate(t *testing.T) {
	type testCase struct {
		name string
		opts []zaplogger.LoggerOpt
	}

	testCases := []testCase{
		{
			name: "default",
			opts: []zaplogger.LoggerOpt{},
		},
		{
			name: "production",
			opts: []zaplogger.LoggerOpt{zaplogger.Production()},
		},
		{
			name: "name_key",
			opts: []zaplogger.LoggerOpt{zaplogger.WithNameKey("test")},
		},
		{
			name: "execution_id_key",
			opts: []zaplogger.LoggerOpt{zaplogger.WithExecutionIDKey("different_key")},
		},
		{
			name: "context_key",
			opts: []zaplogger.LoggerOpt{zaplogger.WithExecutionIDContextKey("context_key")},
		},
		{
			name: "omit_missing_execution_id",
			opts: []zaplogger.LoggerOpt{
				zaplogger.WithOmitExecutionIDWhenMissing(),
				zaplogger.WithExecutionIDContextKey("context_key"),
			},
		},
	}

	for _, tc := range testCases {
		logger, err := zaplogger.NewLogger("test", tc.opts...)
		var _ yall.Logger = logger

		assert.NoError(t, err)
		assert.NotNil(t, logger)

		logger = logger.With("test_name", tc.name)

		// should check the output, but at least thanks to zap DPANIC this would fail if we have a key without a value
		ctx := context.WithValue(context.Background(), yall.ExecutionIDKey, "request_id_test123")
		// logger.Fatal(ctx, "testing execution_id, fatal level")
		// // logger.Fatalnc("testing execution_id, fatal level")
		assert.Panics(t, func() {
			logger.Panic(ctx, "testing execution_id, panic level")
		})
		assert.Panics(t, func() {
			logger.Panicnc("testing execution_id, panic level")
		})
		logger.Error(ctx, "testing execution_id, error level", yall.Error(fmt.Errorf("this is an error")))
		logger.Errornc("testing execution_id, error level", yall.Error(fmt.Errorf("another error")))
		logger.Warn(ctx, "testing execution_id, warn level")
		logger.Warnnc("testing execution_id, warn level")
		logger.Info(ctx, "testing execution_id, info level")
		logger.Infonc("testing execution_id, info level")
		logger.Debug(ctx, "testing execution_id, debug level")
		logger.Debugnc("testing execution_id, debug level")

		// this shouldn't really happen, but better safe than sorry
		//nolint:staticcheck
		logger.Info(nil, "test")
	}
}

func TestExecutionIDFrom(t *testing.T) {
	ctx := context.WithValue(context.Background(), yall.ExecutionIDKey, "request_id_test123")
	logger, _ := zaplogger.NewLogger("test")

	executionID := logger.ExecutionIDFrom(ctx)
	assert.NotEmpty(t, executionID)
	assert.EqualValues(t, "request_id_test123", executionID)
}
