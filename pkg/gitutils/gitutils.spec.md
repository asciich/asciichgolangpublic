# gitutils specifications

This are the specifications for the [`gitutils` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- For functions interacting with files or network there are two implementations, both in a dedicated subpackage. Every function must implemented in both so the functionality stays the same, regardless if executed locally or on a remote machine:
    - [`commandexecutorgit`](commandexecutorgit/README.md):
        - Using `commandexecutor` in combination with shell commands to implement the logic.
        - This also works on other machines over SSH for example.
    - [`nativegit`](nativegit/README.md):
        - Use native golang to handle Git operations.
        - Works only on the local machine.
    - Furthermore, for convenience every function implemented in these packages is as well added to this `gitutils`, calling the [`nativegit`](nativegit/README.md) implementation.
    - For every convenience function add a well commented `Example_<function>_test.go` containing a test for documentation purposes how to use the function
- Generic functions like repository path operations which do not require network or file access are implemented in the subpackage [`gitgeneric`](gitgeneric/README.md).
- For the inmemory representation use appropriate Go types from the standard library.

## Tests

- Test in this directory are used to test all implementations behave the same way. This is why they require a for loop which runs the same test for each implementation.
