# commandexecutorkubernetes package

Implements Kubernetes operations using command executor and kubectl shell commands. Works on remote machines via SSH.

This package is part of the [`kubernetesutils`](../README.md) package hierarchy.

## Implementation Details

This is the **commandexecutor** implementation that uses kubectl commands via the [`commandexecutor`](../../commandexecutor/README.md) interface. This allows execution on remote machines.

For local execution using native Go client libraries, see [`nativekubernetes`](../nativekubernetes/README.md).

## Functions

All functions in this package take an additional `commandExecutor` parameter compared to their [`nativekubernetes`](../nativekubernetes/README.md) counterparts.

## Specifications

For specifications see [kubernetesutils.spec.md](../kubernetesutils.spec.md) and [constitution.md](/constitution.md).
