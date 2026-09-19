# kubeconfigutils

The `kubeconfigutils` package provides utilities for working with Kubernetes kubeconfig files.

## Features

- Load and parse kubeconfig files
- Merge multiple kubeconfig configurations
- Read current kubeconfig respecting `KUBECONFIG` environment variable
- Add additional configurations to existing kubeconfig
- Validate kubeconfig files using kubectl
- Manage contexts, clusters, and users

## Usage

```go
import "gitlab.asciich.ch/tools/asciichgolangpublic.git/pkg/kubernetesutils/kubeconfigutils"
```

### Load a kubeconfig file

```go
ctx := context.Background()
config, err := kubeconfigutils.LoadFromFilePath(ctx, "/path/to/kubeconfig")
```

### Read current kubeconfig

```go
config, err := kubeconfigutils.ReadCurrentKubeConfig(ctx)
```

### Merge configurations

```go
merged, err := kubeconfigutils.MergeConfig(config1, config2)
```

## Specifications

For specifications see [kubeconfigutils.spec.md](kubeconfigutils.spec.md)
