# nativekubernetes package

Non object oriented implementation using the [official client-go library](https://github.com/kubernetes/client-go) to interact with kubernetes.

For the object oriented implementation see [nativekubernetesoo](/pkg/kubernetesutils/nativekubernetesoo/).

## Race Condition Handling

All CRUD operations handle waiting internally to prevent race conditions. For example:

- `CreateNamespace` waits for the default service account to be created
- `DeleteNamespace` waits for the namespace to be fully deleted
- `CreateDeployment` with `WaitForDeploymentAvailable` option waits for all replicas to be ready
- `CreatePod` with `WaitForPodRunning` option waits for the pod to reach Running phase

You do not need to add manual `time.Sleep` calls after these operations.

## Examples

* [Get clientset](Example_GetClientSet_test.go): Get the k8s client-go clientset.
* ConfigMaps:
    * [Create and delete ConfigMap](Example_CreateAndDeleteConfigmap_test.go)
* CronJobs:
    * [Create and delete CronJob](Example_CreateAndDeleteCronjob_test.go)
* Deployments:
    * [Create and delete Deployment](Example_CreateAndDeleteDeployment_test.go)
* Namespaces:
    * [Create and delete namespace](Example_CreateAndDeleteNamespace_test.go)
    * [Get namespace UID](Example_GetNamespaceUID_test.go)
* Pods:
    * [Copy file to pod](Example_CopyFileToPod_test.go): Copy a local file to a pod/container.
    * [Copy local file to pod](Example_CopyLocalFileToPod_test.go): Comprehensive examples showing how to copy files to containers (similar to `kubectl cp`).
        * [Binary file copy](Example_CopyLocalFileToPod_test.go): Copy binary files to containers.
        * [Nested directory copy](Example_CopyLocalFileToPod_test.go): Copy files to nested directories in containers.
    * [Exec](Example_ExecExample_test.go): Run command in already existing pod/container.
        * [Write to stdin of exec command](Example_WriteToStdinOfExecCommand_test.go)
    * [Run command in pod](../Pods_test.go): Run commands in an existing pod using the Pod.RunCommand method.
        * [Native implementation](../Example_RunCommandPod_test.go): Using the native Kubernetes API
        * [Command executor implementation](../Example_RunCommandPod_test.go): Using kubectl exec
        * [Multiple commands](../Example_RunCommandPod_test.go): Running multiple commands in the same pod
        * [Exit code handling](../Example_RunCommandPod_test.go): Checking command exit codes
    * [Get pod logs](Example_GetPodLogs_test.go): Fetch logs from a container running in a pod.
    * [Run temporary pod and get stdout](Example_RunPodAndGetStdout_test.go): How to run a single command in Kubernetes and get it's stdout.
* ReplicaSets:
    * [Create and delete ReplicaSet](Example_CreateAndDeleteReplicaSet_test.go)
* Secrets:
    * [Create and delete Secret](Example_CreateAndDeleteSecret_test.go)
