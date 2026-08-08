# kubernetesutils specifictations

These file contains the specification for the `kubernetesutils` package and it's subpackages.

## Implementation

- For every `kubectl` call use `--context` explicitly to avoid using the wrong cluster by mistake.

## Race condition handling

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
