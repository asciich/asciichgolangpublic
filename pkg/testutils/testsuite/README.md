# testsuite package

Provides a simple way to declare test suites and run them.

## Example Configuration

Here's a complete example showing all available test types including SSH:

```yaml
---
name: "Complete test suite example"
ssh_host: "localhost"                          # Optional: SSH host for remote execution
ssh_user: "root"                               # Optional: SSH user
ssh_port: 22222                                # Optional: SSH port (default: 22)
ssh_skip_host_validation: true                 # Optional: Skip SSH host key validation (for testing)
ssh_private_key_file: "/path/to/private/key"   # Optional: Path to SSH private key file
test_cases:
  # Command test - run any shell command (locally or via SSH if configured)
  - name: "Test echo command"
    test_type: command
    command: echo hello world
    description: "Check if echo command works"
    runbook_links:
      - "https://runbooks.example.com/command-tests"
    hints_for_investigation: "Check if the shell is available and PATH is set correctly."

  # Command test - run on remote SSH server
  - name: "Test curl google"
    test_type: command
    command: curl --fail https://google.com
    description: "Check if we can reach Google using curl"
    runbook_links: "https://runbooks.example.com/network-tests"
    hints_for_investigation: "Verify network connectivity and DNS resolution."

  # Kind cluster exists test - check if a kind cluster exists by name
  - name: "Test kind cluster exists"
    test_type: kind_cluster_exists
    cluster: kind-asciichgolangpublic
    description: "Check if the kind cluster exists"
    runbook_links:
      - "https://runbooks.example.com/kind-cluster"
      - "https://kind.sigs.k8s.io/docs/user/quick-start/"
    hints_for_investigation: "Run 'kind get clusters' to list available clusters."

  # TCP port open test - check if a port is open on a host
  - name: "Test HTTPS port open"
    test_type: tcp_port_open
    port: 443
    host: google.com
    description: "Check if the HTTPS port 443 is open on google.com"
    runbook_links: "https://runbooks.example.com/network-ports"
    hints_for_investigation: "Check firewall rules and network policies."

  # Kubernetes namespace exists test
  - name: "Test default namespace exists"
    test_type: kubernetes_namespace_exists
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if the default namespace exists"
    runbook_links: "https://runbooks.example.com/kubernetes-namespaces"
    hints_for_investigation: "Run 'kubectl get namespaces' to list all namespaces."

  # Kubernetes pod exists test
  - name: "Test pod exists"
    test_type: kubernetes_pod_exists
    resource_name: my-pod
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific pod exists"
    runbook_links: "https://runbooks.example.com/kubernetes-pods"
    hints_for_investigation: "Run 'kubectl get pods -n default' to list pods in the namespace."

  # Kubernetes ReplicaSet exists test
  - name: "Test replicaSet exists"
    test_type: kubernetes_replicaset_exists
    resource_name: my-replicaset
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific ReplicaSet exists"
    runbook_links: "https://runbooks.example.com/kubernetes-replicasets"
    hints_for_investigation: "Run 'kubectl get replicasets -n default' to list ReplicaSets."

  # Kubernetes ConfigMap exists test
  - name: "Test configMap exists"
    test_type: kubernetes_configmap_exists
    resource_name: my-configmap
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific ConfigMap exists"
    runbook_links: "https://runbooks.example.com/kubernetes-configmaps"
    hints_for_investigation: "Run 'kubectl get configmaps -n default' to list ConfigMaps."

  # Kubernetes Secret exists test
  - name: "Test secret exists"
    test_type: kubernetes_secret_exists
    resource_name: my-secret
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific Secret exists"
    runbook_links: "https://runbooks.example.com/kubernetes-secrets"
    hints_for_investigation: "Run 'kubectl get secrets -n default' to list Secrets."

  # Kubernetes Deployment exists test
  - name: "Test deployment exists"
    test_type: kubernetes_deployment_exists
    resource_name: my-deployment
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific Deployment exists"
    runbook_links: "https://runbooks.example.com/kubernetes-deployments"
    hints_for_investigation: "Run 'kubectl get deployments -n default' to list Deployments."

  # Kubernetes CronJob exists test
  - name: "Test cronJob exists"
    test_type: kubernetes_cronjob_exists
    resource_name: my-cronjob
    namespace: default
    cluster: kind-asciichgolangpublic
    description: "Check if a specific CronJob exists"
    runbook_links: "https://runbooks.example.com/kubernetes-cronjobs"
    hints_for_investigation: "Run 'kubectl get cronjobs -n default' to list CronJobs."

  # Kubernetes validate SSH key in secret test - validate that an SSH private key stored in a secret can authenticate
  - name: "Test SSH key validation succeeds"
    test_type: kubernetes_validate_ssh_key_in_secret
    resource_name: ssh-private-key-secret
    secret_key: id_ed25519
    namespace: default
    cluster: kind-asciichgolangpublic
    target_host: ssh-server.default.svc.cluster.local
    target_user: testuser
    target_port: 22
    description: "Check that a valid SSH key in a secret can authenticate"
    runbook_links: "https://runbooks.example.com/kubernetes-ssh-validation"
    hints_for_investigation: "Verify the SSH key in the secret matches the authorized_keys on the target host."
```

## Example Tests

The following example test files demonstrate the usage of the testsuite package:

- `Example_TestGoogleReachable_test.go` - TCP port open and command tests
- `Example_KindClusterExists_test.go` - Kind cluster existence tests
- `Example_KubernetesNamespaceExists_test.go` - Kubernetes namespace existence tests
- `Example_KubernetesPodExists_test.go` - Kubernetes pod existence tests
- `Example_KubernetesReplicaSetExists_test.go` - Kubernetes ReplicaSet existence tests
- `Example_KubernetesConfigMapExists_test.go` - Kubernetes ConfigMap existence tests
- `Example_KubernetesSecretExists_test.go` - Kubernetes Secret existence tests
- `Example_KubernetesDeploymentExists_test.go` - Kubernetes Deployment existence tests
- `Example_KubernetesCronJobExists_test.go` - Kubernetes CronJob existence tests
- `Example_KubernetesValidateSshKeyInSecret_test.go` - Kubernetes SSH key validation in Secret tests

Each example test file contains both localhost and SSH jumphost test scenarios.

## Specifications

For specifications see [testsuite.spec.md](testsuite.spec.md)
