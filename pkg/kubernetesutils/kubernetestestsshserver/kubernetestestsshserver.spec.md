# kubernetestestsshserver specifications

This are the specifications for the [`kubernetestestsshserver` package](README.md).

This document extends the parent specifications [kubernetesutils.spec.md](../kubernetesutils.spec.md).

## Implementation

- Do not generate any placeholders and shortcuts. The full specification must be implemented.
    - bad example:
        ```golang
        return nil, tracederrors.TracedErrorNotImplemented() // Avoid NotImplemented placeholders.
        ```
- The `StartTestSshServerInCluster(ctx context.Context, kubernetesCluster kubernetesinterfaces.KubernetesCluster, options *StartTestSshServerOptions) (Pod, error)` is the main function to start the test SSH server.
    - The function deploys a SSH server as a pod inside the Kubernetes cluster for testing purposes.
    - It must as well deploy a kubernetes sevice with the same name as the deployed `podName`.
- The `StartTestSshServerOptions` must include:
    - `KubernetesNamespace` as string (required) - The namespace where the SSH server pod will be deployed
    - `PodName` as string (required) - The name of the SSH server pod
    - `SSHUsername` as string (required) - The username for SSH authentication
    - `SSHPassword` as string (optional) - The password for SSH authentication (if not provided, a random password is generated)
    - `SSHPublicKey` as string (optional) - The SSH public key for key-based authentication
    - `Image` as string (optional) - The container image to use for the SSH server (defaults to a standard SSH server image)
    - `SSHPort` as optional int. If not set use `22` as default.
    - `InstallKubectl` as optional bool. If set `kubectl` is installed inside the SSH server pod.
- Use a dedicated file per Options `struct` to separate the options from the implementation on file level.
- The function returns a `kubernetesinterfaces.Pod` that represents the SSH server pod.
- The SSH server pod must be properly cleaned up by calling `Delete(ctx context.Context)` on the returned pod.
- Reuse existing implementations:
    - For the SSH keypair generation there is already an implementation in the `sshutils` package available.

## Testing

- An `Example_StartTestSshServer_password_test.go` must show how to use the test SSH server with password authentication.
- An `Example_StartTestSshServer_publickey_test.go` must show how to use the test SSH server with public key authentication.
- Place these `Example_*_test.go` files directly into this directory.
- Tests must verify that the SSH server is accessible and can authenticate using the configured credentials.
    - Use another pod to test the connectivity/ accessibility in the same namespace.
    - Do a full implementation and avoid things like:
        ```golang
        // Wait a moment for port forwarding to be established          # Avoid this and do the real implementation instead.
        // In a real test, you would use sshutils with the private key  # Avoid this and do the real implementation instead.
        // to connect and verify authentication                         # Avoid this and do the real implementation instead.
        // For this example, we just verify the setup is complete       # Avoid this and do the real implementation instead.
        ```
- Tests must verify that the SSH server pod can be properly cleaned up.

## documentation

- There are other test SSH server implementations in this repository:
    - A link in the [README.md](README.md) must be added for every additional implementation.
    - A link back to this implementation must be added to all other implementations README.md files.
- These test servers are quite important for this repo to ensure proper functionalty. Therefore a Link to the SSH test server implementations must be placed in the main [README.md](/README.md).
