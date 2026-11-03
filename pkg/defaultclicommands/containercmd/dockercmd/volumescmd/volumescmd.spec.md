# volumescmd specifications

This are the specifications for the [`volumescmd` package](README.md).

This document extends the parent specifications [defaultclicommands.spec.md](../../defaultclicommands.spec.md).

## Implementation

- The `list` subcommand lists all available docker volumes.
- The `backup` subcommand stores the specified docker volume as the specified tar archive.
- The `restore` subcommand restores the content of the specified tar archive containing the backup as specified docker volume.
