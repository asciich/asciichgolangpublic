# pvccmd

Kubernetes persistent volume claims related commands.

## Usage

```bash
<command> kubernetes pvc
```

## Subcommands

- [`list`](#list) - List all kubernetes persistent volume claims

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

## Specifications

For specifications see [pvccmd.spec.md](pvccmd.spec.md)
