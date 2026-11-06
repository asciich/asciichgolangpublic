# volumescmd specifications

This are the specifications for the [`volumescmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../../defaultclicommands.spec.md).

## Implementation

- The `list` subcommand lists all available docker volumes.
- The `list-files-in-volume` starts a temporary container to list all files in a volume.
    - Include hidden files
    - The output has to be sorted so the list of files becomes deterministic.
