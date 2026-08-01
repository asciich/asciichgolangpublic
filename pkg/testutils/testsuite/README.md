# testsuite package

Provides a simple way to declare test suites and run them.

## Available test types

| test_type                       | Comment                                                    | Example                                                                        |
| ------------------------------- | ---------------------------------------------------------- | ------------------------------------------------------------------------------ | 
| `command`                       | Run an arbitrary command and pass if the exit code is `0`. | [Test Google reachable](./Example_TestGoogleReachable_test.go)                 |
| `tcp_port_open`                 | Passes when the TCP `port` on `host` is open.              | [Test Google reachable](./Example_TestGoogleReachable_test.go)                 |
| `kubernetes_namespace_exists`   | Checks if a Kubernetes namespace exists in a cluster.      | [Test Kubernetes Namespace Exists](./Example_KubernetesNamespaceExists_test.go)|
| `kubernetes_pod_exists`         | Checks if a Kubernetes pod exists in a namespace.          | [Test Kubernetes Pod Exists](./Example_KubernetesPodExists_test.go)            |
| `kubernetes_replicaset_exists`  | Checks if a Kubernetes ReplicaSet exists in a namespace.   | [Test Kubernetes ReplicaSet Exists](./Example_KubernetesReplicaSetExists_test.go)|
| `kubernetes_configmap_exists`   | Checks if a Kubernetes ConfigMap exists in a namespace.    | [Test Kubernetes ConfigMap Exists](./Example_KubernetesConfigMapExists_test.go)|
| `kubernetes_secret_exists`      | Checks if a Kubernetes Secret exists in a namespace.       | [Test Kubernetes Secret Exists](./Example_KubernetesSecretExists_test.go)      |
| `kubernetes_deployment_exists`  | Checks if a Kubernetes Deployment exists in a namespace.   | [Test Kubernetes Deployment Exists](./Example_KubernetesDeploymentExists_test.go)|
| `kubernetes_cronjob_exists`     | Checks if a Kubernetes CronJob exists in a namespace.      | [Test Kubernetes CronJob Exists](./Example_KubernetesCronJobExists_test.go)    |

## Example Configuration

Here's a complete example showing all available test types:

```yaml
---
name: "Complete test suite example"
test_cases:
  # Command test - run any shell command
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

## Field Descriptions

### Common Fields

- **name**: The name of the test case (required)
- **test_type**: The type of test to run (required)
- **description**: A description of what the test does (optional but recommended)

### Test Type Specific Fields

#### command
- **command**: The shell command to execute (required)

#### tcp_port_open
- **host**: The hostname or IP address to check (required)
- **port**: The TCP port number to check (required)

#### kubernetes_namespace_exists
- **namespace**: The name of the Kubernetes namespace to check (required)
- **cluster**: The name of the Kubernetes cluster (required, e.g., `kind-asciichgolangpublic`)

#### kubernetes_pod_exists, kubernetes_replicaset_exists, kubernetes_configmap_exists, kubernetes_secret_exists, kubernetes_deployment_exists, kubernetes_cronjob_exists
- **resource_name**: The name of the Kubernetes resource to check (required)
- **namespace**: The namespace where the resource should exist (required)
- **cluster**: The name of the Kubernetes cluster (required, e.g., `kind-asciichgolangpublic`)
