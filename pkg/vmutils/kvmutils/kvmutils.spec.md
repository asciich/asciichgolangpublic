# kvmutils specifications

This are the specifications for the [`kvmutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Naming:
    - Instead of `poolName` use `storagePoolName` as variable name to be more explicit.
    - Prefer `delete` over `remove`:
        - Use `DeleteVM` instead of `RemoveVM`. This implies as well to use `KvmDeleteVmOptions` instead of `KvmRemoveVmOptions`.
