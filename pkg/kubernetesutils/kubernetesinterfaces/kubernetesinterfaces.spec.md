# Kubernetesinterfaces specifications

This are the specifications for the [`kubernetesinterfaces` package](README.md).

This document extends the parent specifications [kubernetesutils.spec.md](../kubernetesutils.spec.md).

## Implementation

### ConfigMap

- The `ReadYamlValueAsString(ctx context.Context, yamlPath string) (string, error)` is able to directly read a specific value from the `ConfigMap`.
    - The ConfigMap data is converted to a single yaml document (see [kubernetesgeneric](../kubernetesgeneric/README.md)) and the value is read by evaluating the given `yamlPath` as yaml query.
    - Every entry of the ConfigMap data becomes a top level yaml field named like the entry key. Data entries containing valid yaml are embedded as yaml, all other entries as plain string.
    - Since the entry keys become yaml field names, keys containing dots must be quoted in the yaml path: `."values.yaml".image.tag`.
    - An error is returned if the yaml path does not match a value in the ConfigMap.
- The `IsYamlValueAsStringEqual(ctx context.Context, yamlPath string, expectedValue string) (bool, error)` can be used to check if a specific string value matches the `expectedValue`.
    - It reads the value by [ReadYamlValueAsString](#configmap) and compares it to the `expectedValue`.
    - Like `ReadYamlValueAsString` an error is returned if the yaml path does not match a value in the ConfigMap.

## Testing

Testing is done on [`kubernetesutils`](../) level and specified in [kubernetesutils.spec.md](../kubernetesutils.spec.md).
Testing on [`kubernetesutils`](../) level allows to run the test against all available implementations.