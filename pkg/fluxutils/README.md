# Fluxutils

Work with [fluxcd](https://fluxcd.io)

## Subpackages

* [commandexecutorflux](./commandexecutorflux/README.md): Flux operations using command executor and kubectl (works remotely via SSH).
* [nativeflux](./nativeflux/README.md): Native Flux operations using Go libraries (local only).

## Examples

* [Create, delete and watch the flux objects GitRepository, Kustomization and Helmrelease](Example_HandleFluxResources_test.go)
* [Install flux operator in Kubernetes cluster](Example_InstallFluxOperator_test.go)

## Specifications

For specifications see [fluxutils.spec.md](fluxutils.spec.md)