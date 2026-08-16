# kindutils specifications

This are the specifications for the [`kindutils` package](README.md).
        
This document extends the [constitution.md](/constitution.md).

## Implementation

- The cluster and context naming is always a bit confusing and error prone during automation:
    - The clustername itself has no `kind-` prefix, e.g. `my-kind-cluster`.
    - But when looking at `kubectl config get-contexts` both context name and cluster name have the `kind-` prefix, e.g. `kind-my-kind-cluster`.
    - To make the user experience as easy as possible:
        - The `ClusterByNameExists` function must accept both names, with and without prefix. This must be validated by a test.
