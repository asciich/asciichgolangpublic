# kubernetestestsshserver

Deploy a SSH server for testing purposes as pod inside the kubernetes cluster.

There are other test SSH server implementations in this repository:

* [testsshserver](../../sshutils/testsshserver/README.md): Test SSH server running on the local machine.

## Examples

* [Start test SSH server with password authentication](Example_StartTestSshServer_password_test.go): Deploy a SSH server pod for testing with password auth.
* [Start test SSH server with public key authentication](Example_StartTestSshServer_publickey_test.go): Deploy a SSH server pod for testing with key-based auth.

## Specifications

For specifications see [kubernetestestsshserver.spec.md](kubernetestestsshserver.spec.md)
