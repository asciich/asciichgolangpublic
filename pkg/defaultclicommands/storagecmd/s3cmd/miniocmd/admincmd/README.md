# admincmd

Minio admin related commands.

## Specifications

For specifications see [admincmd.spec.md](admincmd.spec.md)

## Commands

### admin

Groups all minio admin related commands.

### admin check-cluster-health

Checks the whole cluster health.

- Logs the cluster status
- Logs all nodes and their status
- Logs all disks and their status
- Shows a success message when everything is ok and exits with code 0
- Shows a fatal error message and exits with non-zero code otherwise
