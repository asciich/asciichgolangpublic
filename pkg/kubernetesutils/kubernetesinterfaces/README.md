# kubernetesinterfaces package

Defines the interfaces for the kubernetes objects like ConfigMap, Namespace, Pod, Deployment and the KubernetesCluster itself.

The implementations of these interfaces are:

* [nativekubernetesoo](../nativekubernetesoo/README.md): Object oriented implementation using the official client-go library (local only).
* [commandexecutorkubernetes](../commandexecutorkubernetes/README.md): Implementation using kubectl commands via the commandexecutor (works remotely via SSH).

This package is part of the [`kubernetesutils`](../README.md) package hierarchy.

## Specifications

For specifications see [kubernetesinterfaces.spec.md](kubernetesinterfaces.spec.md).
