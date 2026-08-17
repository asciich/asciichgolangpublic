# hostutils package

Handle hosts like physical machines and VMs.

While this mostly bases on the `commandexecutorhost` package and therefore hosts are orchestrated using shell commands for the localhost there are two options available:
- `GetLocalCommandExecutorHost()` returns a `commandexecutorhost` based `localhost`
- `GetLocalHost()` returns a `nativehost` based `localhost`.
