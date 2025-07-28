# kubernetesutils

Work with kubernetes.

The implementation using the official k8s client-go is available as:
* [non object oriented implementation](/pkg/kubernetesutils/nativekubernetes/)
* [object oriented implementation](/pkg/kubernetesutils/nativekubernetesoo/) which is on a higher abstraction layer than the non object oriented one.

## Examples

- [ConfigMap by name exists](Example_ConfigmapByNameExists_test.go)
- [List namespace names](Example_ListNamespaceNames_test.go)
- [List node names](Example_ListNodeNames_test.go)
- [Secret by name exists](Example_SecretByNameExists_test.go)
- [Read and write secret](Example_SecretReatAndWrite_test.go)
- [Watch ConfigMap. Get callback on create, update, delete](Example_WatchConfigMap_test.go)
