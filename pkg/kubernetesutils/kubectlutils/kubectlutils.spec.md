# kubectlutils specifications

These file contains the specification for the `kubectlutils` package and it's subpackages.

## Testing

- NO test is allowed to overwrite the system wide `kubectl` or write to `/bin/kubectl` (the default installation path).
- NO test is allowed to call `sudo`.
