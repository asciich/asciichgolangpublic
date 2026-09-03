# fluxutils specifications

This are the specifications for the [`fluxutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- For functions interacting with files or network there are two implementations, both in a dedicated subpackage. Every function must implemented in both so the functionality stays the same, regardless if executed locally or on a remote machine:
    - [`commandexecutorflux`](commandexecutorflux/README.md):
        - Using `commandexecutor` in combination with kubectl and flux shell commands to implement the logic.
        - This also works on other machines over SSH for example.
    - [`nativeflux`](nativeflux/README.md):
        - Use native golang and kubectl client-go to handle Flux operations.
        - Works only on the local machine.
    - Furthermore, for convenience every function implemented in these packages is as well added to this `fluxutils`, calling the [`nativeflux`](nativeflux/README.md) implementation.
    - For every convenience function add a well commented `Example_<function>_test.go` containing a test for documentation purposes how to use the function

## Tests

- Test in this directory are used to test all implementations behave the same way. This is why they require a for loop which runs the same test for each implementation.
