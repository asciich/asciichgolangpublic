# dnsutils specifications

This are the specifications for the [`dnsutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Whenever DNS was resolved log:
    - What was resolved and to which addresses it was resolved.
    - Also include which server resolved it.

## Testing

- In unittests validate the specified logmessages appear correctly.
    - Use `logging.WithLogRecorder(ctx)` to capture the log messages.
