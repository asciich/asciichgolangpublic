# volumescmd

Docker volume related commands.

## Usage

```bash
<command> container docker volumes
```

## Subcommands

- [`list`](#list) - List all available docker volumes
- [`list-files-in-volume`](#list-files-in-volume) - List all files inside a docker volume
- [`backup`](#backup) - Backup a docker volume to a tar archive
- [`restore`](#restore) - Restore a docker volume from a tar archive

### list

List all available docker volumes.

```bash
<command> container docker volumes list
```

### list-files-in-volume

List all files inside a docker volume.

A temporary container is started to list all files inside the given volume.
Hidden files are included and the output is sorted alphabetically.

```bash
<command> container docker volumes list-files-in-volume --volume <volume-name>
```

Usage example:
```bash
<command> container docker volumes list-files-in-volume --volume my-volume
```

### backup

Backup a docker volume to a tar archive.

```bash
<command> container docker volumes backup --volume <volume-name> --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes backup --volume my-volume --archive /backups/my-volume.tar.gz
```

### restore

Restore a docker volume from a tar archive.

```bash
<command> container docker volumes restore --volume <volume-name> --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes restore --volume my-volume --archive /backups/my-volume.tar.gz
```

## Specifications

For specifications see [volumescmd.spec.md](volumescmd.spec.md)

