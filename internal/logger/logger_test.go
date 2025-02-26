package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("level debug", func(t *testing.T) {
		var buf bytes.Buffer

		level := "debug"
		format := "text"

		logger := New(level, format)
		require.NotNil(t, logger)

		logger.logger.Out = &buf
		logger.Debug("this is a debug message")

		output := buf.String()
		require.Contains(t, output, "this is a debug message")
	})

	t.Run("level info", func(t *testing.T) {
		var buf bytes.Buffer

		level := "info"
		format := "text"

		logger := New(level, format)
		require.NotNil(t, logger)

		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this is an info message")

		output := buf.String()
		require.Contains(t, output, "this is an info message")
	})

	t.Run("level warn", func(t *testing.T) {
		var buf bytes.Buffer

		level := "warn"
		format := "text"

		logger := New(level, format)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this should not appear")
		logger.Warn("this is a warning message")

		output := buf.String()
		require.NotContains(t, output, "this should not appear")
		require.Contains(t, output, "this is a warning message")
	})

	t.Run("level error", func(t *testing.T) {
		var buf bytes.Buffer

		level := "error"
		format := "text"

		logger := New(level, format)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Debug("this should not appear")
		logger.Info("this should not appear")
		logger.Warn("this should not appear")
		logger.Error("this is an error message")

		output := buf.String()
		require.NotContains(t, output, "this should not appear")
		require.Contains(t, output, "this is an error message")
	})

	t.Run("invalid level", func(t *testing.T) {
		require.Panics(t, func() {
			level := "invalid"
			format := "text"

			New(level, format)
		}, "expected panic for invalid log level, but none occurred")
	})

	t.Run("json format", func(t *testing.T) {
		var buf bytes.Buffer

		level := "info"
		format := "json"

		logger := New(level, format)
		require.NotNil(t, logger)
		logger.logger.Out = &buf

		logger.Info("this is a info message")

		output := buf.String()
		require.Contains(t, output, `"msg":"this is a info message"`)
	})

	t.Run("invalid format", func(t *testing.T) {
		require.Panics(t, func() {
			level := "info"
			format := "invalid"

			New(level, format)
		}, "expected panic for invalid log format, but none occurred")
	})
}
