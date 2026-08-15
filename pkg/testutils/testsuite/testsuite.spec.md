# testsuite specifications

This are the specifications for the [`testsuite` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Prefer the more specific implementation. Instead of using the more generic `CreateObject` use `Create<TypeName>` like `CreateSecret`.

## Testing

- Being able to run tests on jumphosts is mandatory for this package. Therefore, every `Example_*_test.go` must contain multiple examples:
    - The simple one, testing only on localhost, comes first.
    - The SSH one, demonstrating execution on a jumphost. Refer to `testsuite_ssh_test.go` as a reference for setting up a jumphost-like pod for testing.
        - In both cases, use the `LogRecorder` to capture logs and verify:
            - No SSH commands were used during test execution when testing on localhost.
            - SSH commands were used during test execution when testing on an SSH jumphost.
- For tests requiring a jumphost emulation a SSH server is started inside the kind cluster:
    - The `SetupSSHServerInKind` sets up the SSH server pod and ensures the needed tools are installed.
        - `kubectl` is a mandatory tool and has to be available in the default `PATH` when logged in via SSH.
        - The kubeconfig inside the SSH pod must use the cluster name as the context name to match test expectations.
        - `SetupSSHServerInKind` must verify kubectl installation by running `kubectl version --client` and logging the output.
    - The `TestSetupSSHServerInKind` ensures this procedure works. Make sure it validates the `kubectl` is available and executable in the SSH server pod after the setup as this is needed to run kubernetes related tests.
        - Do not perform the `kubectl` installation in `TestSetupSSHServerInKind`, as it is part of `SetupSSHServerInKind`
        - Validate kubectl by running `kubectl get ns` and verifying the output contains the expected namespace.

## Documentation

- The [README.md](README.md) must contain a section `Example Configuration`:
    - An example entry per available testcase must be present.
    - The testsuite configuration must show all available configuration fields.
- Besides the list of all available `Example_*_test.go` there is no need to document the test setup further in the [README.md](README.md).
    - All `Example_*_test.go` packages must link to the corresponding go file.
