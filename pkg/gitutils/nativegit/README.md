# nativegit package

Implements Git operations using native Go libraries. Works only on the local machine.

This package is part of the [`gitutils`](../README.md) package hierarchy.

## Implementation Details

This is the **native** implementation that uses Go's native git libraries. For remote execution via SSH or other command executors, see [`commandexecutorgit`](../commandexecutorgit/README.md).

## Functions

All functions in this package are also available as convenience functions in the parent [`gitutils`](../README.md) package.

## Specifications

For specifications see [gitutils.spec.md](../gitutils.spec.md) and [constitution.md](/constitution.md).
