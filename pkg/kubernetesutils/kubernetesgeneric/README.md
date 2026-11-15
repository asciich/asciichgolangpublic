# kubernetesgeneric package

Generic helpers to work with kubernetes. The functions in this package operate purely on in-memory data and require no kubernetes connection.

This package is part of the [`kubernetesutils`](../README.md) package hierarchy.

## Functions

* `ConfigMapDataToYamlString(configMapName string, data map[string]string) (string, error)`: Converts the data of a kubernetes ConfigMap to a single yaml document. Every entry of the ConfigMap data becomes a top level yaml field named like the entry key. This way values stored in the ConfigMap can be read by a yaml path like `.my-configmap.hello` afterwards.

## Specifications

For specifications see [kubernetesutils.spec.md](../kubernetesutils.spec.md).
