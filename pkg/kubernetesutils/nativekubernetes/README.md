# nativekubernetes package

Non object oriented implementation using the [official client-go library](https://github.com/kubernetes/client-go) to interact with kubernetes.

For the object oriented implmenetation see [nativekubernetesoo](/pkg/kubernetesutils/nativekubernetesoo/).

## Examples

* [Get clientset](Example_GetClientSet_test.go): Get the k8s client-go clientset.
* [Run temporary pod and get stdout](Example_RunPodAndGetStdout_test.go): How to run a single command in Kubernetes and get it's stdout.
