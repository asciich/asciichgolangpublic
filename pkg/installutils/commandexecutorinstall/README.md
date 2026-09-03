# commandexecutorinstall package

Implements installation operations using command executor and shell commands. Works on remote machines via SSH.

This package is part of the [`installutils`](../README.md) package hierarchy.

## Implementation Details

This is the **commandexecutor** implementation that uses shell commands via the [`commandexecutor`](../../commandexecutor/README.md) interface. This allows execution on remote machines.

For local execution using native Go libraries, a `nativeinstall` package should be created following the patterns in [constitution.md](/constitution.md).

## Functions

All functions in this package take an additional `commandExecutor` parameter compared to their native counterparts.

## Specifications

For specifications see [installutils.spec.md](../installutils.spec.md) and [constitution.md](/constitution.md).
