# commandexecutorhost package

Implements host operations using command executor and shell commands. Works on remote machines via SSH.

This package is part of the [`hostsutils`](../README.md) package hierarchy.

## Implementation Details

This is the **commandexecutor** implementation that uses shell commands via the [`commandexecutor`](../../commandexecutor/README.md) interface. This allows execution on remote machines.

For local execution using native Go libraries, see [`nativehost`](../nativehost/README.md).

## Functions

All functions in this package take an additional `commandExecutor` parameter compared to their [`nativehost`](../nativehost/README.md) counterparts.

## Specifications

For specifications see [hostutils.spec.md](../hostutils.spec.md) and [constitution.md](/constitution.md).
