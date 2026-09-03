# commandexecutorsshclient package

Implements SSH client operations using command executor and shell commands. Works on remote machines via SSH.

This package is part of the [`sshutils`](../README.md) package hierarchy.

## Implementation Details

This is the **commandexecutor** implementation that uses SSH shell commands via the [`commandexecutor`](../../commandexecutor/README.md) interface. This allows execution on remote machines through jump hosts.

For local execution using native Go libraries, see [`nativesshclient`](../nativesshclient/README.md).

## Functions

All functions in this package take an additional `commandExecutor` parameter compared to their [`nativesshclient`](../nativesshclient/README.md) counterparts.

## Specifications

For specifications see [sshutils.spec.md](../sshutils.spec.md) and [constitution.md](/constitution.md).
