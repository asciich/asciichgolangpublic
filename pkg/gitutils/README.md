# gitutils package

Utilities for working with Git repositories, branches, commits, and tags.

## Subpackages

* [commandexecutorgit](./commandexecutorgit/README.md): Git operations using command executor and git shell commands (works remotely via SSH).
* [commandexecutorgitoo](./commandexecutorgitoo/README.md): Object-oriented wrapper around commandexecutorgit.
* [gitgeneric](./gitgeneric/README.md): Generic Git utilities (pure in-memory operations).
* [gitinterfaces](./gitinterfaces/README.md): Git interfaces.
* [gitlabutils](./gitlabutils/README.md): GitLab-specific utilities.
* [nativegit](./nativegit/README.md): Native Git operations using Go libraries (local only).
* [nativegitoo](./nativegitoo/README.md): Object-oriented wrapper around nativegit.

## Implementation Pattern

This package follows the dual-implementation pattern described in [constitution.md](/constitution.md):
- [`nativegit`](nativegit/README.md): Uses native Go git libraries for local execution
- [`commandexecutorgit`](commandexecutorgit/README.md): Uses git shell commands for remote execution

## Specifications

For specifications see [constitution.md](/constitution.md).
