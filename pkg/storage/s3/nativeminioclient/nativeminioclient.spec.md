# nativeminioclient specifications

This are the specifications for the [`nativeminioclient` package](README.md).

This document extends the parent specifications [s3.spec.md](../s3.spec.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- This package provides a native Go implementation for interacting with MinIO S3-compatible storage.
- The implementation uses the MinIO Go client library directly without external `exec` calls.
- All functions are implemented in an idempotent way.
- Functions log "started" and "finished" messages using `logging.LogInfoByCtxf`.

## Usage

- Use the functions in this package for local MinIO client operations.
- For remote operations, use the `commandexecutor` pattern if needed.
