# operateheadscaleapikeyscmd specifications

This are the specifications for the [`operateheadscaleapikeyscmd` package](README.md).

This document extends the [AGENTS.md](/AGENTS.md).

## Implementation

- The `list` CLI command lists the API key prefixes alphabetically ordered.
    - With the `--wide` flag it shows the API key details as table:
        - `ID`
        - `Prefix`
        - `Expiration` (`never` if the key does not expire)
        - `Created`
    - The table is ordered alphabetically by `Prefix`.
