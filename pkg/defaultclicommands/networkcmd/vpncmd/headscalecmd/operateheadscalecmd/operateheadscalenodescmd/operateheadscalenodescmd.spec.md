# operateheadscalenodescmd specifications

This are the specifications for the [`operateheadscalenodescmd` package](README.md).

This document extends the [AGENTS.md](/AGENTS.md).

## Implementation

- The `delete` CLI command can be used to delete a node.
    - It accepts only one argument which is either the node name or node ID.
    - If the argument is a number it is used as node ID, otherwise as node name.
    - Deleting a node which does not exist is not an error (idempotent).
- The `list` CLI command lists the node names alphabetically ordered.
    - With the `--wide` flag it shows the node details as table:
        - `ID`
        - `Hostname`
        - `Name`
        - `User`
        - `Tags`
        - `IP addresses`
        - `Last seen`
        - `Connected`
    - The table is ordered alphabetically by `Name`.
