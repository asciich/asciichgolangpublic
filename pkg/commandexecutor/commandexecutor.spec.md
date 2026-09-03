# commandexecutor specifications

This are the specifications for the [`commandexecutor` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The `commandexecutor` package provides the base interfaces and implementations for executing commands on local and remote systems.
- Multiple implementations are provided:
    - [`commandexecutorgeneric`](commandexecutorgeneric/README.md): Generic command executor functionality.
    - [`commandexecutorexec`](commandexecutorexec/README.md): Execute commands using Go's exec package.
    - [`commandexecutorbash`](commandexecutorbash/README.md): Execute commands using bash shell.
    - [`commandexecutorpowershell`](commandexecutorpowershell/README.md): Execute commands using PowerShell.
- Object-oriented wrappers are provided with the `oo` suffix:
    - [`commandexecutorexecoo`](commandexecutorexecoo/README.md)
    - [`commandexecutorbashoo`](commandexecutorbashoo/README.md)
    - [`commandexecutorpowershelloo`](commandexecutorpowershelloo/README.md)
- The [`commandexecutorinterfaces`](commandexecutorinterfaces/README.md) package defines the core interfaces.

## Security

- Calling exec with unchecked user input leads to security issues.
- To avoid exec calls on the local machine, set the environment variable:
    ```bash
    export ASCIICHGOLANGPUBLIC_AVOID_EXEC=1
    ```

## Tests

- All implementations should be tested to ensure they behave consistently.
- Use the patterns described in [constitution.md](/constitution.md) for cross-implementation testing.
