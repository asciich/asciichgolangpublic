# installkvm

Install KVM related system packages and ensure `libvirtd` is running on localhost.

Both functions require root privileges.

Uses [packagemanager](../../../packagemanager/README.md) to install the packages and [systemdutils](../../../osutils/linuxutils/systemdutils/README.md) to manage the `libvirtd` service.

## Examples

- [Example_Enable_test.go](./Example_Enable_test.go)
- [Example_InstallPackages_test.go](./Example_InstallPackages_test.go)
- [Example_Start_test.go](./Example_Start_test.go)

## Specifications

For specifications see [installkvm.spec.md](./installkvm.spec.md).
