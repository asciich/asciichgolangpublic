# sshutils package

SSH utilities for key generation, management, and client operations.

## Subpackages

* [commandexecutorsshclient](./commandexecutorsshclient/): SSH client using command executor.
* [nativesshclient](./nativesshclient/): Native SSH client implementation.
* [sshoptions](./sshoptions/): SSH configuration options.
* [testsshserver](./testsshserver/): Test SSH server for testing.

## For developers

To run tests use:
```bash
bash -c "cd $(git rev-parse --show-toplevel) && go test -v ./pkg/sshutils/..."
```
