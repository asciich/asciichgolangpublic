package logging_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/asciich/asciichgolangpublic/pkg/contextutils"
	"github.com/asciich/asciichgolangpublic/pkg/logging"
)

func Test_LogRecorder(t *testing.T) {
	t.Run("captures log messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogInfoByCtx(ctx, "operation completed successfully")
		logging.LogInfoByCtxf(ctx, "formatted message: %d", 42)

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "operation completed successfully")
		require.Contains(t, logOutput, "formatted message: 42")
	})

	t.Run("captures multiple messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogInfoByCtx(ctx, "first message")
		logging.LogInfoByCtx(ctx, "second message")
		logging.LogInfoByCtx(ctx, "third message")

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "first message")
		require.Contains(t, logOutput, "second message")
		require.Contains(t, logOutput, "third message")
	})

	t.Run("captures error messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogErrorByCtx(ctx, "an error occurred")
		logging.LogErrorByCtxf(ctx, "error code: %d", 500)

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "an error occurred")
		require.Contains(t, logOutput, "error code: 500")
	})

	t.Run("captures warn messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogWarnByCtxf(ctx, "warning: %s", "low disk space")

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "warning: low disk space")
	})

	t.Run("captures changed messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogChangedByCtx(ctx, "something changed")
		logging.LogChangedByCtxf(ctx, "changed: %s", "value")

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "something changed")
		require.Contains(t, logOutput, "changed: value")
	})

	t.Run("captures good messages", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextVerbose())

		logging.LogGoodByCtx(ctx, "success message")
		logging.LogGoodByCtxf(ctx, "success: %s", "operation completed")

		logOutput := logRecorder.String()
		require.Contains(t, logOutput, "success message")
		require.Contains(t, logOutput, "success: operation completed")
	})

	t.Run("works with nil context", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(nil)
		require.NotNil(t, ctx)
		require.NotNil(t, logRecorder)
	})

	t.Run("empty recorder returns empty string", func(t *testing.T) {
		_, logRecorder := logging.WithLogRecorder(context.Background())
		require.Equal(t, "", logRecorder.String())
	})

	t.Run("does not capture when verbose is disabled", func(t *testing.T) {
		ctx, logRecorder := logging.WithLogRecorder(contextutils.ContextSilent())

		logging.LogInfoByCtx(ctx, "silent message")

		require.Equal(t, "", logRecorder.String())
	})
}
