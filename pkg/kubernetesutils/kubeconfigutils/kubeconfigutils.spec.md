# kubeconfigutils specifications

This are the specifications for the [`kubeconfigutils` package](README.md).

This document extends the parent specifications [kubernetesutils.spec.md](../kubernetesutils.spec.md).

## Implementation

- Must not require the additional config’s context, user, and cluster to share the same name. The context, user, and cluster may each have distinct names (e.g. context `additional-context`, user `additional-user`, and cluster `additional-cluster`), and the merge must resolve them via the context’s references rather than by assuming identical names.
- The `ReadCurrentKubeConfigAsString(ctx context.Context) (string, error)` reads the currently valid kube config like `kubectl` does by respecting the same env var `KUBECONFIG`.
    - Also support the full config directly inside `KUBECONFIG`, not only the path.
- The `ReadCurrentKubeConfig(ctx context.Context) (*KubeConfig, error)` behaves the same way as `ReadCurrentKubeConfigAsString` but with the struct pointer as return value.
- The `AddAdditionalConfigFromString(ctx context.Context, additionalConfig string) error`:
    - Checks if `additionalConfig` is valid.
    - Merges it to the currently valid config (default or set by `KUBECONFIG`):
        - When `KUBECONFIG` itself contains the full config return an error.
    - If no config found and `KUBECONFIG` not set write it as new file to the default location `.kube/config`.

## Testing

- The fact that context, user and cluster are not necessarily equal requires a test case for each combination.
