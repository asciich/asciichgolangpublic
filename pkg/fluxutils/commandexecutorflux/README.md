# commandexecutorflux package

Implements Flux operations using command executor and shell commands. Works on remote machines via SSH.

This package is part of the [`fluxutils`](../README.md) package hierarchy.

## Implementation Details

This is the **commandexecutor** implementation that uses kubectl and flux shell commands via the [`commandexecutor`](../../commandexecutor/README.md) interface. This allows execution on remote machines.

For local execution using native Go libraries, see [`nativeflux`](../nativeflux/README.md).

## Functions

All functions in this package take an additional `commandExecutor` parameter compared to their [`nativeflux`](../nativeflux/README.md) counterparts.

## Specifications

For specifications see [fluxutils.spec.md](../fluxutils.spec.md) and [constitution.md](/constitution.md).
