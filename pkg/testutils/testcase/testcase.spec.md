# testcase specifications

This are the specifications for the [`testcase` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- Every testcase is defined in it's own `*.go` file.
- The `Run` function of every testcase implementation receives the `commandExecutor *commandexecutorinterfaces.CommandExecutor` to use:
    - If it's local host the implementation must use tha native, go only implementations.
    - Otherwise the `commandExecutor` has to be used to execute the test setps.
    - This way we ensure all test cases can run on a remote host as well.
