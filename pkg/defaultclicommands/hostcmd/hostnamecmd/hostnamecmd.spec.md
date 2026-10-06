# hostnamecmd specifications

This are the specifications for the [`hostnamecmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../defaultclicommands.spec.md).

## Implementation

- At least these subcommands need to be implemented:
    - `get-hostname` to get the current hostname.
    - `set-hostname` to set the hostname persistently.
- The commands must wire up with the implementation in the [`nativehost`](../../hostsutils/nativehost/README.md) package.
    - Use `nativehost.NewNativeHost()` to get the host instance for localhost operations.
    - For `get-hostname`: Retrieve the current hostname using `hostnamectl --static hostname` command execution.
    - For `set-hostname`: Use the `SetHostName()` method which executes `sudo hostnamectl set-hostname` persistently.
- Each cobra.Command must include a usage example in the `Long` description as specified in [defaultclicommands.spec.md](../defaultclicommands.spec.md).
- The `set-hostname` command must:
    - Accept exactly one argument (the new hostname).
    - Validate that the hostname is not empty.
    - Be idempotent - skip if hostname is already set to the requested value.
    - Use `sudo` for the hostnamectl set-hostname command only if requested by the `--sudo` parameter.
- The `get-hostname` command must:
    - Output the current static hostname to stdout.
    - Not require any arguments.
