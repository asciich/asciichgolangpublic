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
* [Create and delete Role](Example_CreateAndDeleteRole_test.go)
* [Create and delete ClusterRole](Example_CreateAndDeleteClusterRole_test.go)
* [Run single command in temporary pod](Example_RunSingleCommandPod_test.go)
* [Run single command in temporary pod with secret](Example_RunSingleCommandPodWithSecret_test.go)
* [Run single command in temporary pod with secret as file](Example_RunSingleCommandPodWithSecretAsFile_test.go)
* [Run command in temporary pod](Example_RunCommandInTemporaryPod_test.go)
* [Copy file to pod](Example_CopyFileToPod_test.go): Copy a local file to a container running in a pod (similar to `kubectl cp`).
    * [Native Kubernetes implementation](Example_CopyFileToPod_test.go): Using the native Kubernetes API
    * [Command Executor implementation](Example_CopyFileToPod_test.go): Using kubectl cp command
    * [Nested directory copy](Example_CopyFileToPod_test.go): Copy files to nested directories
* [Copy file from pod](Example_CopyFileFromPod_test.go): Copy a file from a container running in a pod to the local filesystem.
    * [Native Kubernetes implementation](Example_CopyFileFromPod_test.go): Using the native Kubernetes API
    * [Command Executor implementation](Example_CopyFileFromPod_test.go): Using kubectl cp command
    * [Round-trip example](Example_CopyFileFromPod_test.go): Copy file to pod and back to verify integrity
* [Run command in existing pod](Pods_test.go): Execute commands in a running pod using the Pod.RunCommand method.
    * [Native Kubernetes implementation](Example_RunCommandPod_test.go): Using the native Kubernetes API
    * [Command Executor implementation](Example_RunCommandPod_test.go): Using kubectl exec command
    * [Multiple commands example](Example_RunCommandPod_test.go): Running multiple commands in the same pod
    * [Exit code example](Example_RunCommandPod_test.go): Checking command exit codes
* The examples to exec/ run commands as additional process inside a container are in [nativekubernetes](./nativekubernetes/README.md)
* Namespaces:
    * [Create and delete namespace](nativekubernetes/Example_CreateAndDeleteNamespace_test.go)
    * [Get namespace UID](nativekubernetes/Example_GetNamespaceUID_test.go)
    * [List namespace names](Example_ListNamespaceNames_test.go)
* [List node names](Example_ListNodeNames_test.go)
* [Secret by name exists](Example_SecretByNameExists_test.go)
* [Read and write secret](Example_SecretReadAndWrite_test.go)
* [Validate SSH key in secret](Example_ValidateSSHKeyInSecret_test.go): Test if a Kubernetes secret contains a valid SSH private key by attempting to SSH into a target host.
* [Watch ConfigMap. Get callback on create, update, delete](Example_WatchConfigMap_test.go)
* [Wait for pod ready](Example_WaitPodReady_test.go)

## Specifications

See [kubernetesutils.spec.md](kubernetesutils.spec.md)
