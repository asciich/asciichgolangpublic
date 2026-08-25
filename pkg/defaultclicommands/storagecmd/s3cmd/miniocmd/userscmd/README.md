# userscmd

Users related commands for MinIO.

## Specifications

For specifications see [userscmd.spec.md](userscmd.spec.md)

## Commands

### `users create`

Create a new user.

```
Usage:
    <command> storage s3 minio users create <username>
```

Flags:
- `--password` - Password for the new user (required)
- `--read-only` - Give the user read only permissions
- `--keep-current-password-if-user-exists` - Keep the current password if the user already exists

### `users delete`

Delete a user.

```
Usage:
    <command> storage s3 minio users delete <username>
```

### `users list`

List all users.

```
Usage:
    <command> storage s3 minio users list
```
