# localgitoo package

Legacy object-oriented implementation of a locally available git repository using the built-in git functionality of the [go-git](https://github.com/go-git/go-git) library.

This package is part of the [`gitutils`](../README.md) package hierarchy.

`LocalGitRepository` implements the [`gitinterfaces.GitRepository`](../gitinterfaces/README.md) interface for local repositories. For new code prefer the [`nativegitoo`](../nativegitoo/README.md) package.

`GitRepositoriesService` provides convenience functions to clone repositories into local directories.

## Specifications

For specifications see [gitutils.spec.md](../gitutils.spec.md) and [constitution.md](/constitution.md).
