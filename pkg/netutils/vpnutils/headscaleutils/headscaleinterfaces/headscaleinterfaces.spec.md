# headscaleinterfaces specifications

This are the specifications for the [`headscaleinterfaces` package](README.md).

This document extends the [AGENTS.md](/AGENTS.md).

## Implementation

- The `Node` interface must define getters for the fields:
    - `id`: `number`
    - `hostname`: `string`
    - `name`: `string`
    - `machineKey`: `string`
    - `nodeKey`: `string`
    - `user`: `User | string`
    - `tags`: `string[]`
    - `ipAddresses`: `string[]`
    - `ephemeral`: `boolean`
    - `lastSeen`: `Date | null`
    - `expiry`: `Date | null` (`0001-01-01 00:00:00` maps to `null`)
    - `online`: `boolean`
    - `expired`: `boolean`
- The `User` interface must define getters for the fields:
    - `id`: `number`
    - `name`: `string` (display name, can be empty)
    - `username`: `string`
    - `email`: `string` (can be empty)
    - `created`: `Date`
- The `ApiKey` interface must define getters for the fields:
    - `id`: `number`
    - `prefix`: `string`
    - `expiration`: `Date | null` (`0001-01-01 00:00:00` maps to `null`)
    - `created`: `Date`
- The `Headscale` interface must define:
    - `ListNodeNames(ctx context.Context) ([]string, error)` to retrieve only the node names.
    - `ListNodes(ctx context.Context) ([]Node, error)` to retrieve the nodes with details.
    - `ListUserNames(ctx context.Context) ([]string, error)` to retrieve only the usernames.
    - `ListUsers(ctx context.Context) ([]User, error)` to retrieve the users with details.
    - `ListApiKeys(ctx context.Context) ([]ApiKey, error)` to retrieve the API keys with details.
    - `DeleteNodeByName(ctx context.Context, nodeName string) error` to delete a node by its name.
    - `DeleteNodeById(ctx context.Context, id uint64) error` to delete a node by its ID.
