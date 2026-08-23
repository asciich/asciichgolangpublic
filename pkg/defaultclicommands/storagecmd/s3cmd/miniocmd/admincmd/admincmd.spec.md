# minio admincmd specifications

## Implementation

- The commands in this package only orchestrate the logic of the `nativeminioclient` client.
    - If a function is missing implement it in the `nativeminioclient` and only reuse it here.
- The command `admin` is used to group all minio admin related commands.
    - The `check-cluster-health` must is used to check the whole cluster health.
        - Log the cluster status
        - Log all nodes and its status
        - Log all disks and its status
        - Show a LogGood... message when everything is ok and exit 0.
        - Use LogFatal with an error message otherwise (this automatically exists != 0).
