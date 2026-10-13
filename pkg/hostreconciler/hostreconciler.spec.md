# hostreconciler specifications

## Implementation

- The reconcile logic is implemented in `ReconcileOnce`.
- `ReconcileOnce` always performs these steps:
    - Reads the desired hostname and applies it.
