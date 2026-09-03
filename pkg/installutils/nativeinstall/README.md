# nativeinstall package

Implements installation operations using native Go libraries. Works only on the local machine.

This package is part of the [`installutils`](../README.md) package hierarchy.

## Implementation Details

This is the **native** implementation that uses Go's native libraries for file and HTTP operations. For remote execution via command executors, see [`commandexecutorinstall`](../commandexecutorinstall/README.md).

## Functions

All functions in this package are also available as convenience functions in the parent [`installutils`](../README.md) package.

## Specifications

For specifications see [installutils.spec.md](../installutils.spec.md) and [constitution.md](/constitution.md).
