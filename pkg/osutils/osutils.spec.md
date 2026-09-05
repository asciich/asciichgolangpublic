# osutils specifications

## Implementation:

- The `TurnSwapOff(ctx context.Context)` convenience function calls `commandexecutorlinuxuserutils.TurnSwapOff(ctx context.Context, commandexecutorexecoo.Exec())` to:
    - Turn off swap on the system if needed.
