# commandexecutoruserutils specifications

This are the specifications for the [`commandexecutoruserutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- `func IsRunningAsRoot(ctx context.Context, commandexecutorinterfaces.CommandExecutor) (bool, error)` can be used to check if running as root.

## Testing

- Every exported function must be tested.
- Use the `dockerutils` package to start temporary docker containers to test the implementation. Do not run tests against the host itself.
