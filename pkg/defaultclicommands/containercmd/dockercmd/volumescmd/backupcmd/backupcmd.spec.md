# backupcmd specifications

This are the specifications for the [`backupcmd` package](README.md).

This document extends the parent specifications [volumescmd.spec.md](../volumescmd.spec.md).

## Implementation

- The `backup` subcommand stores the specified docker volume as the specified tar archive.
- The `backup-multi` subcommand stores the specified docker volumes as specified tar archive.
    - The archive contains all files of all backed up volumes in a way it can simply restore all of them.
    - The `--volume` can be specified multiple times to list all volumes to back up.
    - The `--all` parameter can be used to backup all docker volumes into a single file.
- The `restore` subcommand restores the content of the specified tar archive containing the backup as specified docker volume.
- The `restore-multi` subcoammdn restores all volumes in the specified tar archive.
    - If the tar archive does not contain the information where to restore an error is trown.
