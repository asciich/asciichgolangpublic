# logging Specifications

## Implementation

- The `LogRecorder` captures log messages in memory while still emitting them normally. This is useful for testing scenarios where you need to assert on log output.
    - Example usage:
        ```go
        ctx, logRecorder := logging.WithLogRecorder(ctx)
        logging.LogInfoByCtx(ctx, "operation completed successfully")
        require.Contains(t, logRecorder.String(), "operation completed successfully")
        ```
    - Implement unittests for it.
