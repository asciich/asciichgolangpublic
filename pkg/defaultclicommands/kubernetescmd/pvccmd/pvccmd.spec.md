# pvccmd specifications

This are the specifications for the [`pvccmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../../defaultclicommands.spec.md).

## Implementation

- The `pvc` command must provide this subcommand:
    - The `list` subcommand lists all PVCs.
        - Default is for all namespaces unless `--namespace` is used to select a specific namespace.
- The `list-files-in-pvc` starts a temporary container to list all files in a pvc.
    - Include hidden files
    - The output has to be sorted so the list of files becomes deterministic.
