# volumescmd

Docker volume related commands.

## Usage

```bash
<command> container docker volumes
```

## Subcommands

- [`list`](#list) - List all available docker volumes
- [`list-files-in-volume`](#list-files-in-volume) - List all files inside a docker volume
- [`backup`](#backup) - Backup and restore docker volumes

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

Backup and restore docker volumes. See the [`backupcmd` package](./backupcmd/README.md) for all subcommands.

```bash
<command> container docker volumes backup
```

## Specifications

For specifications see [volumescmd.spec.md](volumescmd.spec.md)

