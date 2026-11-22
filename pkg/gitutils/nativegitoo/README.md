# nativegitoo package

Object-oriented wrapper around [`nativegit`](../nativegit/README.md) package.

This package is part of the [`gitutils`](../README.md) package hierarchy.

## Implementation Details

This package provides an object-oriented API for Git operations. It wraps the procedural functions from [`nativegit`](../nativegit/README.md) in structs with methods.

`LocalGitRemote` represents a remote of a local git repository and wraps the `git.Remote` type of the [go-git](https://github.com/go-git/go-git) library.

## Specifications

For specifications see [gitutils.spec.md](../gitutils.spec.md) and [constitution.md](/constitution.md).
