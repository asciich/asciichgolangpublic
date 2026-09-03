# hostutils package

Handle hosts like physical machines and VMs.

While this mostly bases on the `commandexecutorhost` package and therefore hosts are orchestrated using shell commands for the localhost there are two options available:
- `GetLocalCommandExecutorHost()` returns a `commandexecutorhost` based `localhost`
- `GetLocalHost()` returns a [`nativehost`](nativehost/README.md) based `localhost`.

## Subpackages

* [commandexecutorhost](./commandexecutorhost/README.md): Host operations using command executor (works remotely via SSH).
* [nativehost](./nativehost/README.md): Native host implementation using Go standard library.

## Specifications

For specifications see [hostsutils.spec.md](hostsutils.spec.md)
