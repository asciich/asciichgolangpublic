# commandexecutorinterfaces package

Defines the interfaces for command executor implementations.

This package is part of the [`commandexecutor`](../README.md) package hierarchy.

## Implementation Details

This package defines the `CommandExecutor` interface that is implemented by various command executor backends like:
- [`commandexecutorgeneric`](../commandexecutorgeneric/README.md)
- [`commandexecutorexec`](../commandexecutorexec/README.md)
- [`commandexecutorbash`](../commandexecutorbash/README.md)
- [`commandexecutorpowershell`](../commandexecutorpowershell/README.md)

## Interfaces

The main interface is `CommandExecutor` which provides methods for executing commands on local or remote systems.
