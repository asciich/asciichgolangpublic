# hostreconciler specifications

This are the specifications for the [`hostreconciler` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The reconcile logic is implemented in `ReconcileOnce`.
- `ReconcileOnce` performs these steps:
    - Reads the desired hostname and applies it.
    
