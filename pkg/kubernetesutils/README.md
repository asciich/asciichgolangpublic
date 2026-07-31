# kubernetesutils

Work with kubernetes.

The implementation using the official k8s client-go is available as:
* [non object oriented implementation](/pkg/kubernetesutils/nativekubernetes/)
* [object oriented implementation](/pkg/kubernetesutils/nativekubernetesoo/) which is on a higher abstraction layer than the non object oriented one.
The implementation using exec to call `kubectl` or other commands is useful when jumphosts are in between. It is available as:
* [commandexecutorkubernetes](/pkg/kubernetesutils/commandexecutorkubernetes/)

## Examples

* [Check ConfigMap by name exists](Example_ConfigmapByNameExists_test.go)
* [Check CronJob by name exists](Example_CheckCronJobByNameExists_test.go)
* [CronJob by name exists](Example_CronjobByNameExists_test.go)
* [Check Secret by name exists](Example_CheckSecretByNameExists_test.go)
* [Check Namespace by name exists](Example_CheckNamespaceByNameExists_test.go)
* [Check Pod by name exists](Example_CheckPodByNameExists_test.go)
* [Check ReplicaSet by name exists](Example_CheckReplicaSetByNameExists_test.go)
* [Check Deployment by name exists](Example_CheckDeploymentByNameExists_test.go)
* [Create and delete ConfigMap](Example_CreateAndDeleteConfigMap_test.go)
* [Create and delete CronJob](Example_CreateAndDeleteCronJob_test.go)
* [Create and delete Deployment](Example_CreateAndDeleteDeployment_test.go)
* [Create and delete Pod](Example_CreateAndDeletePod_test.go)
* [Create and delete ReplicaSet](Example_CreateAndDeleteReplicaSet_test.go)
* [Run single command in temporary pod](Example_RunSingleCommandPod_test.go)
* [Run single command in temporary pod with secret](Example_RunSingleCommandPodWithSecret_test.go)
* [Run single command in temporary pod with secret as file](Example_RunSingleCommandPodWithSecretAsFile_test.go)
* The examples to exec/ run commands as additional process inside a container are in [nativekubernetes](./nativekubernetes/README.md)
* Namespaces:
    * [Create and delete namespace](nativekubernetes/Example_CreateAndDeleteNamespace_test.go)
    * [Get namespace UID](nativekubernetes/Example_GetNamespaceUID_test.go)
    * [List namespace names](Example_ListNamespaceNames_test.go)
* [List node names](Example_ListNodeNames_test.go)
* [Secret by name exists](Example_SecretByNameExists_test.go)
* [Read and write secret](Example_SecretReadAndWrite_test.go)
* [Watch ConfigMap. Get callback on create, update, delete](Example_WatchConfigMap_test.go)
