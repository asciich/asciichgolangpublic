# volumescmd

Docker volume related commands.

## Usage

```bash
<command> container docker volumes
```

## Subcommands

- [`list`](#list) - List all available docker volumes
- [`backup`](#backup) - Backup a docker volume to a tar archive
- [`restore`](#restore) - Restore a docker volume from a tar archive

### list

List all available docker volumes.

```bash
<command> container docker volumes list
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

