# kindutils specifications

This are the specifications for the [`kindutils` package](README.md).
        
This document extends the [constitution.md](/constitution.md).

## Implementation

- The cluster and context naming is always a bit confusing and error prone during automation:
    - The clustername itself has no `kind-` prefix, e.g. `my-kind-cluster`.
    - But when looking at `kubectl config get-contexts` both context name and cluster name have the `kind-` prefix, e.g. `kind-my-kind-cluster`.
    - To make the user experience as easy as possible:
        - The `ClusterByNameExists` function must accept both names, with and without prefix. This must be validated by a test.
- The kind binary version and the node image version have to be compatible with each other:
    - The kind binary is pinned by `InstallKind`.
    - The default node image `DefaultNodeImage` matches this pinned kind binary version.
    - `CreateCluster` uses `DefaultNodeImage` unless a different `NodeImage` is set in the `CreateClusterOptions`.
    - This avoids `kind create cluster` failures like `could not find a log line that matches "Reached target .*Multi-User System.*|detected cgroup v1"` caused by a kind binary that is older than the latest `kindest/node` default image.
