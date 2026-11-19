# backupcmd

Docker volume backup and restore commands.

## Usage

```bash
<command> container docker volumes backup
```

## Subcommands

- [`create`](#create) - Backup a docker volume to a tar archive
- [`create-multi`](#create-multi) - Backup multiple docker volumes into one tar archive
- [`restore`](#restore) - Restore a docker volume from a tar archive
- [`restore-multi`](#restore-multi) - Restore all volumes contained in a multi volume tar archive

### create

Backup a docker volume to a tar archive.

```bash
<command> container docker volumes backup create --volume <volume-name> --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes backup create --volume my-volume --archive /backups/my-volume.tar.gz
```

### create-multi

Backup multiple docker volumes into one tar archive. The archive contains a manifest.json and the content of every volume. Use `restore-multi` to restore all volumes of the archive at once.

The `--volume` flag can be specified multiple times. Use `--all` to backup all docker volumes.

```bash
<command> container docker volumes backup create-multi (--volume <volume-name>... | --all) --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes backup create-multi --volume vol-a --volume vol-b --archive /backups/volumes.tar.gz
<command> container docker volumes backup create-multi --all --archive /backups/all-volumes.tar.gz
```

### restore

Restore a docker volume from a tar archive.

```bash
<command> container docker volumes backup restore --volume <volume-name> --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes backup restore --volume my-volume --archive /backups/my-volume.tar.gz
```

### restore-multi

Restore all docker volumes contained in a multi volume tar archive. The volumes to restore are taken from the manifest.json inside the archive.

```bash
<command> container docker volumes backup restore-multi --archive <archive-path>
```

Usage example:
```bash
<command> container docker volumes backup restore-multi --archive /backups/volumes.tar.gz
```

## Specifications

For specifications see [backupcmd.spec.md](backupcmd.spec.md)


