# testsuite specifications

This are the specifications for the [`testsuite` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Prefer the more specific implementation. Instead of using the more generic `CreateObject` use `Create<TypeName>` like `CreateSecret`.

### TestCases

- Each test case is defined in the YAML configuration under `test_cases`.
- Every test case must have the following common fields:
    - `name`: A unique, descriptive name for the test case.
    - `test_type`: The type of test to execute.
    - `description`: A human-readable description of what the test validates.
- If SSH configuration is provided at the suite level (`ssh_host`, `ssh_user`, etc.), command-based tests execute on the remote host via SSH.

#### Available Test Types

##### `command`
- Executes a shell command locally or via SSH (if SSH is configured).
- Required fields:
    - `command`: The shell command to execute.
- The test passes if the command exits with code 0.

##### `kind_cluster_exists`
- Checks whether a kind (Kubernetes IN Docker) cluster exists by name.
- Required fields:
    - `cluster`: The name of the kind cluster to check.
- The test passes if the kind cluster exists (i.e., appears in `kind get clusters` output).

##### `kubernetes_namespace_exists`
- Checks whether a Kubernetes namespace exists in the specified cluster.
- Required fields:
    - `namespace`: The namespace name to check.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the namespace exists.

##### `kubernetes_pod_exists`
- Checks whether a Kubernetes pod exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The pod name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the pod exists.

##### `kubernetes_replicaset_exists`
- Checks whether a Kubernetes ReplicaSet exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The ReplicaSet name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the ReplicaSet exists.

##### `kubernetes_configmap_exists`
- Checks whether a Kubernetes ConfigMap exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The ConfigMap name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the ConfigMap exists.

##### `kubernetes_secret_exists`
- Checks whether a Kubernetes Secret exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The Secret name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the Secret exists.

##### `kubernetes_deployment_exists`
- Checks whether a Kubernetes Deployment exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The Deployment name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the Deployment exists.

##### `kubernetes_cronjob_exists`
- Checks whether a Kubernetes CronJob exists in the specified namespace and cluster.
- Required fields:
    - `resource_name`: The CronJob name.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
- The test passes if the CronJob exists.

##### `kubernetes_validate_ssh_key_in_secret`
- Validates that an SSH private key stored in a Kubernetes Secret can successfully authenticate against a target host.
- Required fields:
    - `resource_name`: The Secret name containing the SSH private key.
    - `secret_key`: The key within the Secret data that holds the private key.
    - `namespace`: The namespace to look in.
    - `cluster`: The Kubernetes cluster context name.
    - `target_host`: The SSH host to authenticate against.
    - `target_user`: The SSH user to authenticate as.
    - `target_port`: The SSH port on the target host.
- The test passes if SSH authentication succeeds using the key from the Secret.

##### `tcp_port_open`
- Checks whether a TCP port is reachable on a given host.
- Required fields:
    - `host`: The target hostname or IP address.
    - `port`: The TCP port number to check.
- The test passes if a TCP connection can be established.

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
    - `KinD` related tests can't run over SSH as the pod in the cluster does not see the docker socket nor anything else to check for `KinD` functionality.

## Documentation

- The [README.md](README.md) must contain a section `Example Configuration`:
    - An example entry per available testcase must be present.
    - The testsuite configuration must show all available configuration fields.
- Besides the list of all available `Example_*_test.go` there is no need to document the test setup further in the [README.md](README.md).
    - All `Example_*_test.go` packages must link to the corresponding go file.
