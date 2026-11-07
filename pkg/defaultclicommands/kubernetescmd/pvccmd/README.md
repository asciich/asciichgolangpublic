# pvccmd

Kubernetes persistent volume claims related commands.

## Usage

```bash
<command> kubernetes pvc
```

## Subcommands

- [`list`](#list) - List all kubernetes persistent volume claims
- [`list-files-in-pvc`](#list-files-in-pvc) - List all files inside a kubernetes persistent volume claim

### list

List all kubernetes persistent volume claims.

Default is to list the persistent volume claims of all namespaces. Use `--namespace` to only list the persistent volume claims of a specific namespace.

```bash
<command> kubernetes pvc list
<command> kubernetes pvc list --namespace <namespace>
```

Usage example:
```bash
<command> kubernetes pvc list --namespace longhorn-system
```

### list-files-in-pvc

List all files inside a kubernetes persistent volume claim.

A temporary pod is started to list all files inside the given persistent volume claim.
Hidden files are included and the output is sorted alphabetically.

```bash
<command> kubernetes pvc list-files-in-pvc --pvc <pvc-name> [--namespace <namespace>]
```

Usage example:
```bash
<command> kubernetes pvc list-files-in-pvc --pvc my-pvc --namespace longhorn-system
```

## Specifications

For specifications see [pvccmd.spec.md](pvccmd.spec.md)
