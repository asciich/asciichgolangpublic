# admincmd specifications

This are the specifications for the [`admincmd` package](README.md).

This document extends the [constitution.md](/constitution.md).

## Implementation

- The commands in this package only orchestrate the logic of the [`nativeminioclient`](../../../../storage/s3/nativeminioclient/README.md) client.
    - If a function is missing implement it in the [`nativeminioclient`](../../../../storage/s3/nativeminioclient/README.md) and only reuse it here.
- The command `admin` is used to group all minio admin related commands.
    - The `check-cluster-health` must is used to check the whole cluster health.
        - Log the cluster status
        - Log all nodes and its status.
            - The uptime of the node must be in a human readable format. Reuse existing libraries in this repo.
        - Log all disks and its status
        - Show a LogGood... message when everything is ok and exit 0.
        - Use LogFatal with an error message otherwise (this automatically exists != 0).
