# tlscmd specifications

## Implementation

- This package does not implement logic itself:
    - This package is used to wire up the cobra.Cobra CLI commands with the actual implementation.
    - If additional functionality is required implement it in the `tlsutils` or the correct sub package of `tlsutils`.
- The `tls` command must provide these subcommands:
    - `get-from-server` to get the certificate directly from the `--hostname` on given `--port`:
        - Use the default TLS/SSL port 443 if the `--port` is not specified.
    - `show-info` to show the info from a local file.
        - If `-` is given read from stdin instead.
