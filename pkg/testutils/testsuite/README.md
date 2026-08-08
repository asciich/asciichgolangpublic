# testsuite package

Provides a simple way to declare test suites and run them.

## Example Configuration

Here's a complete example showing all available test types including SSH:

```yaml
---
name: "Complete test suite example"
ssh_host: "localhost"      # Optional: SSH host for remote execution
ssh_user: "root"           # Optional: SSH user
test_cases:
  # Command test - run any shell command (locally or via SSH if configured)
  - name: "Test echo command"
    test_type: command
    command: echo hello world
    description: "Check if echo command works"

  # Command test - run on remote SSH server
  - name: "Test curl google"
    test_type: command
    command: curl --fail https://google.com
    description: "Check if we can reach Google using curl"

  # TCP port open test - check if a port is open on a host
  - name: "Test HTTPS port open"
    test_type: tcp_port_open
    port: 443
    host: google.com
    description: "Check if the HTTPS port 443 is open on google.com"

  # Kubernetes namespace exists test
  - name: "Test default namespace exists"
    test_type: kubernetes_namespace_exists
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if the default namespace exists"

  # Kubernetes pod exists test
  - name: "Test pod exists"
    test_type: kubernetes_pod_exists
    resource_name: my-pod
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific pod exists"

  # Kubernetes ReplicaSet exists test
  - name: "Test replicaSet exists"
    test_type: kubernetes_replicaset_exists
    resource_name: my-replicaset
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific ReplicaSet exists"

  # Kubernetes ConfigMap exists test
  - name: "Test configMap exists"
    test_type: kubernetes_configmap_exists
    resource_name: my-configmap
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific ConfigMap exists"

  # Kubernetes Secret exists test
  - name: "Test secret exists"
    test_type: kubernetes_secret_exists
    resource_name: my-secret
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific Secret exists"

  # Kubernetes Deployment exists test
  - name: "Test deployment exists"
    test_type: kubernetes_deployment_exists
    resource_name: my-deployment
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific Deployment exists"

  # Kubernetes CronJob exists test
  - name: "Test cronJob exists"
    test_type: kubernetes_cronjob_exists
    resource_name: my-cronjob
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific CronJob exists"
```

## Example Tests

The following example test files demonstrate the usage of the testsuite package:

- `Example_TestGoogleReachable_test.go` - TCP port open and command tests
- `Example_KubernetesNamespaceExists_test.go` - Kubernetes namespace existence tests
- `Example_KubernetesPodExists_test.go` - Kubernetes pod existence tests
- `Example_KubernetesReplicaSetExists_test.go` - Kubernetes ReplicaSet existence tests
- `Example_KubernetesConfigMapExists_test.go` - Kubernetes ConfigMap existence tests
- `Example_KubernetesSecretExists_test.go` - Kubernetes Secret existence tests
- `Example_KubernetesDeploymentExists_test.go` - Kubernetes Deployment existence tests
- `Example_KubernetesCronJobExists_test.go` - Kubernetes CronJob existence tests

Each example test file contains both localhost and SSH jumphost test scenarios.

## Specification

See [testsuite.spec.md](testsuite.spec.md)

## Specifications

For specifications see [testsuite.spec.md](testsuite.spec.md)
