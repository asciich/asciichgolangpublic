# commandexecutorfile specifications

This are the specifications for the [`commandexecutorfile` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Do not use `bash` for evaluation, use `sh` instead since available on more systems.
- For getting the size of a file do NOT read the while file to avoid performance issues. Always use commands getting the file size from the filesystem instead.

## Testing

- A unit test for all functions performing a `RunCommand` must be added.
- For tests requiring `sudo` use docker container as provided by the `dockerutils` package.
- Test all functions with the `alpine:latest` image. It uses busybox and this package must support it.
