# kubectlutils specifications

This are the specifications for the [`kubectlutils` package](README.md).

This document extends the parent specifications [kubernetesutils.spec.md](../kubernetesutils.spec.md).

## Testing

- NO test is allowed to overwrite the system wide `kubectl` or write to `/bin/kubectl` (the default installation path).
- NO test is allowed to call `sudo`.
