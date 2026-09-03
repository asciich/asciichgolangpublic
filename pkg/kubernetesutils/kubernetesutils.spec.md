# kubernetesutils specifications

This are the specifications for the [`kubernetesutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- For every `kubectl` call use `--context` explicitly to avoid using the wrong cluster by mistake.
- The `Namespace` interface must provide a `Delete(context.Context) error` function to delete itself.
- All the functions must be implemented in an idempotent way so there is no need to add a lot of if/else handing:
    - Bad example:
        ```golang
        // Check if namespace exists, create if not  // Not needed since Create itself is already idempotent.
        exists, err := namespace.Exists(ctx)         // Not needed since Create itself is already idempotent.
        if err != nil {                              // Not needed since Create itself is already idempotent.
            return nil, err                          // Not needed since Create itself is already idempotent.
        }                                            // Not needed since Create itself is already idempotent.
        if !exists {                                 // Not needed since Create itself is already idempotent.
            err = namespace.Create(ctx)              // This function must be already idempotent.
            if err != nil {                          // As a consequence the err != nil check is needed only once.
                return nil, err                      // As a consequence the err != nil check is needed only once.
            }                                        // As a consequence the err != nil check is needed only once.
        }                                            // Not needed since Create itself is already idempotent.
        ```
- Be very verbose in your logmessages:
    - Always include the object, the namespace and the kubernetes cluster name:
        - Instead of:
            ```
            Wait until pod 'test-ssh-client-password' in namespace 'test-ssh-password' is deleted finished.
            ```
        - Print
            ```
            Wait until pod 'test-ssh-client-password' in namespace 'test-ssh-password' of kubernetes cluster 'clustername' is deleted finished.
            ```

### Race condition handling

- CRUD operations on kubernetes objects must not lead to race conditions.
    - The need for sleeps after a CRUD operation call must be avoided:
        Instead of:
            ```golang
            // Create namespace
            _, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
            require.NoError(t, err)

            // Wait until default sa in created namespace exists.
            time.Sleep(10 * time.Second)

            // Go on with your actual logic
            ```
        the waiting must be handled in the `CreateNAmespaceByName`:
            ```golang
            // Create namespace
            _, err := kubernetes.CreateNamespaceByName(ctx, namespaceName)
            require.NoError(t, err)

            // Go on with your actual logic, no "manual" wait needed.
            ```

## Testing

- Every mentioned implementation requires a test for both the `nativekubernetes` and `commandexecutorkubernetes` based implementation.
- If a test requires a SSH server do not use the `testsshserver` package, start a test SSH server as pod in the cluster using `kubernetestestsshserver`.
- Preferably use the shared cluster for testing:
    - Run the initailization of it in every test case using the shared cluster as one of the first steps:
        ```golang
        // -----
        // Prepare test environment start ...
        kubernetesCluster, err := kindutils.GetOrCreateSharedCluster(ctx)
        if err != nil {
            require.NoError(t, err)
        }
        // ... prepare test environment finished.
        // -----
        ```
