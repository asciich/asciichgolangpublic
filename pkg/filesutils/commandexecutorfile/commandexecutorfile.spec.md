# commandexecutorfile specifications

This are the specifications for the [`commandexecutorfile` package](README.md).

This document extends the parent specifications [filesutils.spec.md](../filesutils.spec.md).

## Implementation

- Do not use `bash` for evaluation, use `sh` instead since available on more systems.

## Testing

- A unit test for all functions performing a `RunCommand` must be added.
