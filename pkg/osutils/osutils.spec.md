# osutils specifications

This are the specifications for the [`osutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The `TurnSwapOff(ctx context.Context)` convenience function calls `linuxutils.TurnSwapOff(ctx, commandexecutorexecoo.Exec())` to:
    - Turn off swap on the system if needed.
