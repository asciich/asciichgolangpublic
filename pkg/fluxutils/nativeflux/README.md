# nativeflux package

Implements Flux operations using native Go libraries. Works only on the local machine.

This package is part of the [`fluxutils`](../README.md) package hierarchy.

## Implementation Details

This is the **native** implementation that uses Go's native libraries and kubectl for Flux operations. For remote execution via command executors, see [`commandexecutorflux`](../commandexecutorflux/README.md).

## Functions

All functions in this package are also available as convenience functions in the parent [`fluxutils`](../README.md) package.

## Specifications

For specifications see [fluxutils.spec.md](../fluxutils.spec.md) and [constitution.md](/constitution.md).
